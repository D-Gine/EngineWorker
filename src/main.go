package main

import (
	ece "github.com/D-Gine/ECE"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	registry := ece.NewRegistry()

	router.Use(func(c *gin.Context) {
		c.Set("registry", registry)
		c.Next()
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	router.Run()
}
