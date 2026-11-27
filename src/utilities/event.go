package utilities

import "fmt"

// Event for the Entity Component Event to be queued, emitted, subscribed to...
type APIEvent struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`
}

func (ev APIEvent) String() string {
	return ev.Name + " " + fmt.Sprintf("%v", ev.Metadata)
}
