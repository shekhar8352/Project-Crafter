package main

import (
	"backend/database"
	"os"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var crafterCollection *mongo.Collection = database.OpenCollection(database.Client, "crafter")

func main() {
	// Start the server
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	router := gin.New()

	router.Run(":" + port)
}
