package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"

	// "backend/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// Create Gin router
	r := gin.Default()

	// Register routes
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

