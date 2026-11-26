package eventqueue

import (
	"EngineWorker/src/engine"
	"EngineWorker/src/utilities"
	"net/http"
	"slices"
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
	eng := c.MustGet("engine").(*engine.Engine)

	var event utilities.APIEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eng.Do(func(r *ece.Registry) {
		r.EnqueueEventByName(event.Name, event)
	})
	c.JSON(http.StatusAccepted, gin.H{"event": event.Name, "queued": true})
}

func Subscribe(c *gin.Context) {
	eng := c.MustGet("engine").(*engine.Engine)
	var subscriber Subscriber
	if err := c.ShouldBindJSON(&subscriber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callback := func(r *ece.Registry, event ece.Event) any {
		utilities.Load_scripts(r, event, subscriber.Script)
		return nil
	}

	eng.Do(func(r *ece.Registry) {
		r.SubscribeToEventByName(subscriber.Event, callback, subscriber.Priority)
	})

	c.JSON(http.StatusCreated, gin.H{"event": subscriber.Event, "script": subscriber.Script, "priority": subscriber.Priority})
}

func Process(c *gin.Context) {
	eng := c.MustGet("engine").(*engine.Engine)

	event_nb, err := strconv.Atoi(c.DefaultQuery("n", "-1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of events to process"})
		return
	}

	eng.DoAsync(func(r *ece.Registry) {
		if event_nb > 0 {
			r.ProcessManyEvents(event_nb)
		} else {
			r.ProcessEvents()
		}
	})

	c.JSON(http.StatusAccepted, gin.H{"status": "processing"})
}

func EnqueueBreakpoint(c *gin.Context) {
	eng := c.MustGet("engine").(*engine.Engine)

	eng.Do(func(r *ece.Registry) {
		r.EnqueueEvent(ece.BreakpointEvent{})
	})

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

func GetEventQueue(c *gin.Context) {
	eng := c.MustGet("engine").(*engine.Engine)

	events := engine.Ask(eng, func(r *ece.Registry) []ece.Pair[string, ece.Event] {
		return slices.Clone(r.GetEventQueue())
	})

	answer := gin.H{"queue": make([]map[string]any, 0)}
	for _, event := range events {
		if apiEvent, ok := event.Second.(utilities.APIEvent); ok {
			answer["queue"] = append(answer["queue"].([]map[string]any), map[string]any{
				"Name":     apiEvent.Name,
				"Metadata": apiEvent.Metadata,
			})
		} else if _, ok := event.Second.(ece.BreakpointEvent); ok {
			answer["queue"] = append(answer["queue"].([]map[string]any), map[string]any{
				"Name":     "BreakpointEvent",
				"Metadata": make(map[string]any),
			})
		} else {
			continue
		}
	}
	c.JSON(http.StatusOK, answer)
}
