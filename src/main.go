package main

import (
	"EngineWorker/src/controllers/eventqueue"

	ece "github.com/D-Gine/ECE"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	registry := ece.NewRegistry()

	router.Use(func(c *gin.Context) {
		c.Set("registry", registry)
		c.Next()
	})

	events := router.Group("/events")
	{
		events.GET("/process", eventqueue.Process)
		events.POST("/enqueue", eventqueue.Enqueue)
		events.POST("/subscribe", eventqueue.Subscribe)
	}

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	router.Run()
}
