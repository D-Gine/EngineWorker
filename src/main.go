package main

import (
	"os"

	ece "github.com/D-Gine/ECE"
	"github.com/D-Gine/EngineWorker/src/database"
	"github.com/D-Gine/EngineWorker/src/debug"
	"github.com/gin-gonic/gin"
)

func main() {
	port := ":" + os.Getenv("API_PORT")

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	registry := ece.NewRegistry()

	// Database setup
	debug.PrintMsg("=== DATABASE ===")
	err := database.ConnectDatabase()
	if err != nil {
		debug.PrintErrorMsg("failed to connect to database:" + err.Error())
	}
	database.RetreiveRoomInformations()

	router.Use(func(c *gin.Context) {
		c.Set("registry", registry)
		c.Next()
	})
	debug.PrintMsg("")

	if err := router.Run(port); err != nil {
		debug.PrintErrorMsg("Failed to start server:" + err.Error())
		os.Exit(1)
	}
}
