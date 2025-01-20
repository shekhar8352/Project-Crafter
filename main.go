package main

import (
	"crafter/database"
	"fmt"
	"os"

	"crafter/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var crafterCollection *mongo.Collection = database.OpenCollection(database.Client, "crafter")

func main() {
	// Start the server
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("PORT is not found in the environment variable")
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Run(":" + port)
}
