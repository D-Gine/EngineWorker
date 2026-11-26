package utilities

import (
	"EngineWorker/src/debug"
	"fmt"
	"time"

	ece "github.com/D-Gine/ECE"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

const (
	maxExecutionSteps = 100_000
	scriptTimeout     = 5 * time.Second
)

// EmitBuiltin is a builtin to emit an event from a starlark script
func EmitBuiltin(reg *ece.Registry) *starlark.Builtin {
	return starlark.NewBuiltin("emit", func(
		thread *starlark.Thread,
		b *starlark.Builtin,
		args starlark.Tuple,
		kwargs []starlark.Tuple,
	) (starlark.Value, error) {
		var name string

		if err := starlark.UnpackPositionalArgs("event", args, nil, 1, &name); err != nil {
			return starlark.MakeInt(-1), err
		}

		event := &APIEvent{Name: name, Metadata: map[string]any{}}
		event.Name = name

		for _, kv := range kwargs {
			v, err := StarlarkToGo(kv[1])
			if err == nil {
				event.Metadata[string(kv[0].(starlark.String))] = v
			}
		}
		reg.EmitEventByName(event.Name, *event)
		return starlark.MakeInt(0), nil
	})
}

// Load_scripts executes a starlark script that can access all the attributes of an event given in parameter
func Load_scripts(reg *ece.Registry, event ece.Event, script string) {
	thread := &starlark.Thread{Name: "Script thread"}
	thread.SetMaxExecutionSteps(maxExecutionSteps)

	timer := time.AfterFunc(scriptTimeout, func() {
		thread.Cancel("script execution timed out")
	})
	defer timer.Stop()

	ev, ok := event.(APIEvent)

	predeclared := starlark.StringDict{
		"emit": EmitBuiltin(reg),
	}
	if ok {
		predeclared["event_name"] = starlark.String(ev.Name)

		for k, v := range ev.Metadata {
			sv, err := GoToStarlark(v)
			if err != nil {
				debug.Error(err)
				return
			}
			predeclared[k] = sv
		}
	} else {
		println("event is not an APIEvent:", fmt.Sprintf("%T", event))
	}

	_, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, "script.star", script, predeclared)
	if err != nil {
		debug.Error(err)
		return
	}
}
