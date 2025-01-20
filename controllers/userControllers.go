package controllers

import (
	"crafter/database"
	"crafter/models"
	"crafter/utils"
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

func returnResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": "Operation successful",
		"data":    data,
	})
}

// returnError sends a JSON response with the provided status code and error message.
func returnError(c *gin.Context, statusCode int, errMessage string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": errMessage,
	})
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Bind the JSON body to the user model
		if err := c.BindJSON(&user); err != nil {
			returnError(c, http.StatusBadRequest, err.Error())
			return
		}

		// Validate the user data
		validationErr := validate.Struct(user)
		if validationErr != nil {
			returnError(c, http.StatusBadRequest, validationErr.Error())
			return
		}

		// Check if the email already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			returnError(c, http.StatusInternalServerError, "error occurred while checking for the email")
			return
		}

		if count > 0 {
			returnError(c, http.StatusBadRequest, "this email already exists")
			return
		}

		// Ensure password is not null
		if user.Password == nil {
			returnError(c, http.StatusBadRequest, "password cannot be null")
			return
		}

		// Hash the password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Set creation and update timestamps
		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		// Generate a new ObjectID for the user
		user.ID = primitive.NewObjectID()
		User_id := user.ID.Hex()

		// Generate JWT tokens
		token, refreshToken, err := utils.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, User_id)
		if err != nil {
			returnError(c, http.StatusInternalServerError, "error generating tokens")
			return
		}

		// Set tokens in the user model
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// Insert the user into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			returnError(c, http.StatusInternalServerError, "user item was not created")
			return
		}

		// Return success response
		returnResponse(c, http.StatusOK, gin.H{
			"msg":    "user item was created successfully",
			"result": resultInsertionNumber,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			returnError(c, http.StatusBadRequest, err.Error())
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			fmt.Println(err)
			returnError(c, http.StatusUnauthorized, "Cannot find user with these credentials")
			return
		}

		if foundUser.Password == nil || user.Password == nil {
			returnError(c, http.StatusBadRequest, "Password cannot be null")
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			returnError(c, http.StatusBadRequest, msg)
			return
		}

		User_id := foundUser.ID.Hex()

		token, refreshToken, err := utils.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, User_id)
		if err != nil {
			returnError(c, http.StatusInternalServerError, "error generating tokens")
			return
		}

		err = utils.UpdateAllTokens(token, refreshToken, User_id)
		if err != nil {
			returnError(c, http.StatusInternalServerError, "error updating tokens")
			return
		}

		returnResponse(c, http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage := 10
		page := 1

		// Get 'recordPerPage' query parameter if provided
		if rpp, err := strconv.Atoi(c.Query("recordPerPage")); err == nil && rpp > 0 {
			recordPerPage = rpp
		}

		// Get 'page' query parameter if provided
		if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
			page = p
		}

		// Calculate the starting index for pagination
		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{"$match", bson.D{}}}
		countStage := bson.D{{"$count", "total_count"}}
		paginationStages := []bson.D{
			{{"$skip", startIndex}},
			{{"$limit", recordPerPage}},
		}

		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"first_name", 1},
				{"last_name", 1},
				{"email", 1},
				{"user_type", 1},
				{"experience_level", 1},
				{"date_of_birth", 1},
				{"resume_urls", 1},
				{"college", 1},
				{"current_company", 1},
			}},
		}

		// Fetch total user count
		totalUsersCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, countStage})
		if err != nil {
			returnError(c, http.StatusInternalServerError, "error occurred while fetching user count")
			return
		}

		var totalCount []bson.M
		if err := totalUsersCursor.All(ctx, &totalCount); err != nil {
			returnError(c, http.StatusInternalServerError, "error occurred while counting users")
			return
		}

		totalUserCount := 0
		if len(totalCount) > 0 {
			if count, ok := totalCount[0]["total_count"].(int32); ok {
				totalUserCount = int(count)
			}
		}

		// Fetch paginated list of users
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, paginationStages[0], paginationStages[1], projectStage,
		})

		if err != nil {
			returnError(c, http.StatusInternalServerError, "error occurred while listing user items")
			return
		}

		var users []bson.M
		if err := result.All(ctx, &users); err != nil {
			returnError(c, http.StatusInternalServerError, "error fetching users")
			return
		}

		// Return success response
		returnResponse(c, http.StatusOK, gin.H{
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
			returnError(c, http.StatusBadRequest, "user_id parameter is required")
			return
		}

		// Convert the userId to a MongoDB ObjectID
		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			returnError(c, http.StatusBadRequest, "Invalid ObjectID")
			return
		}

		// Fetch the user from the database
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				returnError(c, http.StatusNotFound, "user not found")
			} else {
				returnError(c, http.StatusInternalServerError, "error occurred while retrieving user")
			}
			return
		}

		// Check if the user data is valid
		if user.ID.IsZero() {
			returnError(c, http.StatusInternalServerError, "retrieved user data is invalid")
			return
		}

		// Return the user object on success
		returnResponse(c, http.StatusOK, user)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fmt.Println("Invalid ObjectID string:", err)
			return
		}

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.ID = primitive.NilObjectID // Assuming the field is named ID in your User model

		updateData := bson.M{
			"$set": bson.M{
				"first_name":       user.First_name, // Add only fields that need to be updated
				"last_name":        user.Last_name,
				"email":            user.Email,
				"date_of_birth":    user.Date_of_birth,
				"user_type":        user.UserType,
				"experience_level": user.Experience,
				"college":          user.College,
				"current_company":  user.Current_company,
			},
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user_id}, updateData)
		if err != nil {
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
