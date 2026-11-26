package engine

import (
	"EngineWorker/src/debug"
	"fmt"

	ece "github.com/D-Gine/ECE"
)

type Engine struct {
	tasks chan func(*ece.Registry)
}

func runTask(registry *ece.Registry, task func(*ece.Registry)) {
	defer func() {
		if r := recover(); r != nil {
			debug.ErrorMessage(fmt.Sprintf("ECE task panicked: %v", r))
		}
	}()
	task(registry)
}

func NewEngine(registry *ece.Registry) *Engine {
	e := &Engine{
		tasks: make(chan func(*ece.Registry), 64),
	}
	go func() {
		for task := range e.tasks {
			runTask(registry, task)
		}
	}()
	return e
}

func (e *Engine) Do(fn func(*ece.Registry)) {
	done := make(chan struct{})
	e.tasks <- func(r *ece.Registry) {
		defer close(done)
		fn(r)
	}
	<-done
}

func (e *Engine) DoAsync(fn func(*ece.Registry)) {
	e.tasks <- fn
}

func Ask[T any](e *Engine, fn func(*ece.Registry) T) T {
	res := make(chan T, 1)
	e.tasks <- func(r *ece.Registry) {
		defer close(res)
		res <- fn(r)
	}
	return <-res
}
