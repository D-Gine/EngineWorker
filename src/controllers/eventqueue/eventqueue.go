package eventqueue

import (
	"EngineWorker/src/utilities"
	"net/http"
	"strconv"

	ece "github.com/D-Gine/ECE"
	"github.com/gin-gonic/gin"
)

type Subscriber struct {
	Event    string `json:"event"`
	Script   string `json:"script"`
	Priority int    `json:"priority"`
}

func Enqueue(c *gin.Context) {
	registry := c.MustGet("registry").(*ece.Registry)

	var event utilities.APIEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registry.EnqueueEventByName(event.Name, event)
	c.JSON(http.StatusAccepted, gin.H{"event": event.Name, "queued": true})
}

func Subscribe(c *gin.Context) {
	registry := c.MustGet("registry").(*ece.Registry)
	var subscriber Subscriber
	if err := c.ShouldBindJSON(&subscriber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callback := func(r *ece.Registry, event ece.Event) any {
		utilities.Load_scripts(r, event, subscriber.Script)
		return nil
	}

	registry.SubscribeToEventByName(subscriber.Event, callback, subscriber.Priority)

	c.JSON(http.StatusCreated, gin.H{"event": subscriber.Event, "script": subscriber.Script, "priority": subscriber.Priority})
}

func Process(c *gin.Context) {
	registry := c.MustGet("registry").(*ece.Registry)

	event_nb, err := strconv.Atoi(c.DefaultQuery("n", "-1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of events to process"})
		return
	}
	var processed []ece.Pair[string, ece.Event]
	if event_nb > 0 {
		processed = registry.ProcessManyEvents(event_nb)
	} else {
		processed = registry.ProcessEvents()
	}

	if len(processed) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No events to process"})
		return
	}

	answer := gin.H{"Processed events": make([]string, 0)}
	for _, event := range processed {
		answer["Processed events"] = append(answer["Processed events"].([]string), event.First)
	}
	c.JSON(http.StatusOK, answer)
}
