package utilities

import (
	"fmt"

	ece "github.com/D-Gine/ECE"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func Load_scripts(reg *ece.Registry, event ece.Event, script string) {
	thread := &starlark.Thread{Name: "Script thread"}

	ev, ok := event.(APIEvent)

	predeclared := starlark.StringDict{}
	if ok {
		predeclared["event_name"] = starlark.String(ev.Name)

		for k, v := range ev.Metadata {
			sv, err := GoToStarlark(v)
			if err != nil {
				panic(err)
			}
			predeclared[k] = sv
		}
	} else {
		println("event is not an APIEvent:", fmt.Sprintf("%T", event))
	}

	_, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, "script.star", script, predeclared)
	if err != nil {
		panic(err)
	}
}
