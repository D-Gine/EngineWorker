package main

import (
	"EngineWorker/src/controllers/eventqueue"
	"os"

	"EngineWorker/src/database"
	"EngineWorker/src/debug"

	ece "github.com/D-Gine/ECE"
	"github.com/gin-gonic/gin"
)

func main() {
	port := ":" + os.Getenv("API_PORT")

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	registry := ece.NewRegistry()

	// Database setup
	debug.Message("=== DATABASE ===")
	err := database.ConnectDatabase()
	if err != nil {
		debug.ErrorMessage("failed to connect to database:" + err.Error())
	}
	database.RetreiveRoomInformations()

	router.Use(func(c *gin.Context) {
		c.Set("registry", registry)
		c.Next()
	})
	debug.Message("")

	events := router.Group("/events")
	{
		events.GET("/process", eventqueue.Process)
		events.GET("/queue", eventqueue.GetEventQueue)

		events.POST("/enqueue", eventqueue.Enqueue)
		events.POST("/subscribe", eventqueue.Subscribe)

		events.PUT("/breakpoint", eventqueue.EnqueueBreakpoint)
	}

	if err := router.Run(port); err != nil {
		debug.ErrorMessage("Failed to start server:" + err.Error())
		os.Exit(1)
	}
}
