package utilities

import "fmt"

type APIEvent struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`
}

func (ev APIEvent) String() string {
	return ev.Name + " " + fmt.Sprintf("%v", ev.Metadata)
}
