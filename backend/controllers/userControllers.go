package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	if userPassword == "" || providedPassword == "" {
		return false, "password cannot be null"
	}

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Password is incorrect"
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email already exists"})
			return
		}

		if user.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be null"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, err := utils.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating tokens"})
			return
		}
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "user item was created successfully", "result": resultInsertionNumber})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if foundUser.Password == nil || user.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be null"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating tokens"})
			return
		}

		err = utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating tokens"})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage := 10
		page := 1

		if rpp, err := strconv.Atoi(c.Query("recordPerPage")); err == nil && rpp > 0 {
			recordPerPage = rpp
		}
		if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
			page = p
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{"$match", bson.D{}}}
		countStage := bson.D{{"$count", "total_count"}}
		paginationStages := []bson.D{
			{{"$skip", startIndex}},
			{{"$limit", recordPerPage}},
		}

		projectStage := bson.D{
			{"$project", bson.D{
				{"user_id", 1},
				{"first_name", 1},
				{"last_name", 1},
				{"email", 1},
				{"user_type", 1},
				{"experience_level", 1},
				{"date_of_birth", 1},
				{"resume_urls", 1},
				{"college", 1},         // Include College
				{"current_company", 1}, // Include Current_company // Include Refresh_Token
			}},
		}

		totalUsersCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, countStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching user count"})
			return
		}
		var totalCount []bson.M
		if err := totalUsersCursor.All(ctx, &totalCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while counting users"})
			return
		}
		totalUserCount := 0
		if len(totalCount) > 0 {
			if count, ok := totalCount[0]["total_count"].(int32); ok {
				totalUserCount = int(count)
			}
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, paginationStages[0], paginationStages[1], projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}

		var users []bson.M
		if err := result.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total_count":   totalUserCount,
			"users":         users,
			"page":          page,
			"recordPerPage": recordPerPage,
		})
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
			return
		}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving user"})
			}
			return
		}

		if user.ID.IsZero() {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "retrieved user data is invalid"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		fmt.Println(userId)
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
			return
		}

		var user models.User	
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": user})
		if err != nil {
			fmt.Println(err)
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating user"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
	}
}
