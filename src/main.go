package main

import (
	"fmt"
	"os"

	ece "github.com/D-Gine/ECE"
	"github.com/D-Gine/EngineWorker/src/database"
	"github.com/gin-gonic/gin"
)

func main() {
	port := ":" + os.Getenv("API_PORT")

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	registry := ece.NewRegistry()

	// Database setup
	err := database.ConnectDatabase()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}

	router.Use(func(c *gin.Context) {
		c.Set("registry", registry)
		c.Next()
	})

	if err := router.Run(port); err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
