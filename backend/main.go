package main

import (
	"os"
	"github.com/gin-gonic/gin"
)

func main() {
	// Start the server
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	router := gin.New()


	router.Run(":" + port)
}
