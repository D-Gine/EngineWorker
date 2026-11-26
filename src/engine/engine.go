package engine

import (
	"EngineWorker/src/debug"
	"fmt"

	ece "github.com/D-Gine/ECE"
)

// Engine serializes all operations on an ece.Registry in a dedicated goroutine.
type Engine struct {
	tasks chan func(*ece.Registry)
}

// runTask executes a single task and recovers from panics so the worker loop keeps running.
func runTask(registry *ece.Registry, task func(*ece.Registry)) {
	defer func() {
		if r := recover(); r != nil {
			debug.ErrorMessage(fmt.Sprintf("ECE task panicked: %v", r))
		}
	}()
	task(registry)
}

// NewEngine starts the background ECE goroutine and returns an Engine handle.
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

// Do runs fn on the ECE goroutine and waits for completion.
func (e *Engine) Do(fn func(*ece.Registry)) {
	done := make(chan struct{})
	e.tasks <- func(r *ece.Registry) {
		defer close(done)
		fn(r)
	}
	<-done
}

// DoAsync queues fn on the ECE goroutine without waiting for completion.
func (e *Engine) DoAsync(fn func(*ece.Registry)) {
	e.tasks <- fn
}

// Ask runs fn on the ECE goroutine and returns its result.
func Ask[T any](e *Engine, fn func(*ece.Registry) T) T {
	res := make(chan T, 1)
	e.tasks <- func(r *ece.Registry) {
		defer close(res)
		res <- fn(r)
	}
	return <-res
}
