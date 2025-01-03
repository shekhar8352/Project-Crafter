package main

import (
	"backend/models"
	"backend/controllers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mockUserCollection = []models.User{}

func stringPointer(s string) *string {
	return &s
}

// TestSignUp_Success tests the successful creation of a new user
//
// The test creates a new user with all the required fields and
// asserts that the response status code is 200 and the response
// body contains the message "user item was created successfully".
func TestSignUp_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp())

	user := models.User{
		First_name:    stringPointer("John"),
		Last_name:     stringPointer("Doe"),
		Date_of_birth: stringPointer("02-09-2002"),
		Password:      stringPointer("Password123"),
		Email:         stringPointer("john.doe@example.com"),
		UserType:      models.Professional,
		Experience:    models.Fresher,
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
		ID:            primitive.NewObjectID(),
		User_id:       primitive.NewObjectID().Hex(),
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user item was created successfully")
}


// TestSignUp_MissingFields tests the signup endpoint with missing required fields
//
// The test attempts to create a new user without the required first name field
// and expects the response to be a 400 Bad Request with an error message.
func TestSignUp_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp())

	user := models.User{
		Last_name:     stringPointer("Doe"),
		Date_of_birth: stringPointer("02-09-2002"),
		Password:      stringPointer("Password123"),
		Email:         stringPointer("john.doe@example.com"),
		UserType:      models.Professional,
		Experience:    models.Fresher,
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}


// TestSignUp_InvalidEnum tests the signup endpoint with invalid UserType and ExperienceLevel enums.
//
// The test attempts to create a new user with invalid UserType and ExperienceLevel enums
// and expects the response to be a 400 Bad Request with an error message.
func TestSignUp_InvalidEnum(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp())

	user := models.User{
		First_name:    stringPointer("John"),
		Last_name:     stringPointer("Doe"),
		Date_of_birth: stringPointer("02-09-2002"),
		Password:      stringPointer("Password123"),
		Email:         stringPointer("john.doe@example.com"),
		UserType:      "InvalidType",       // Invalid UserType
		Experience:    "InvalidExperience", // Invalid ExperienceLevel
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}


// TestSignUp_DuplicateEmail tests the signup endpoint with a duplicate email
//
// The test assumes there's already a user in the system with the same email and
// attempts to create a new user with the same email. The test expects the response
// to be a 400 Bad Request with an error message.
func TestSignUp_DuplicateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Assume there's already a user in the system
	existingUser := models.User{
		Email: stringPointer("john.doe@example.com"),
	}
	mockUserCollection = append(mockUserCollection, existingUser)

	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp())

	user := models.User{
		First_name:    stringPointer("John"),
		Last_name:     stringPointer("Doe"),
		Date_of_birth: stringPointer("02-09-2002"),
		Password:      stringPointer("Password123"),
		Email:         stringPointer("john.doe@example.com"), // Duplicate email
		UserType:      models.Professional,
		Experience:    models.Fresher,
	}

	jsonValue, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Error marshalling user: %v", err)
	}
	req, err := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "this email already exists")
}
