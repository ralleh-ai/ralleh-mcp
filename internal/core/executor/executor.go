package executor

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
)

// Task is one bounded unit of source work.
type Task[T any] struct {
	SourceID string
	Run      func(context.Context) (T, error)
}

// Outcome captures the result or failure for a task. A failed source must not
// fail the whole request; callers return partials with diagnostics.
type Outcome[T any] struct {
	SourceID string
	Value    T
	Err      error
	Duration time.Duration
}

// RunBounded executes tasks with hard global/per-source deadlines and bounded
// concurrency. It never launches more than b.MaxConcurrency workers.
func RunBounded[T any](ctx context.Context, b budget.Budget, tasks []Task[T]) []Outcome[T] {
	if b.MaxSources > 0 && len(tasks) > b.MaxSources {
		tasks = tasks[:b.MaxSources]
	}
	if len(tasks) == 0 {
		return nil
	}
	globalCtx, cancel := context.WithTimeout(ctx, b.GlobalTimeout)
	defer cancel()

	concurrency := b.MaxConcurrency
	if concurrency <= 0 || concurrency > len(tasks) {
		concurrency = len(tasks)
	}
	jobs := make(chan Task[T])
	outcomes := make(chan Outcome[T], len(tasks))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for task := range jobs {
			start := time.Now()
			sourceCtx, sourceCancel := context.WithTimeout(globalCtx, b.PerSourceTimeout)
			value, err := task.Run(sourceCtx)
			sourceCancel()
			if errors.Is(globalCtx.Err(), context.DeadlineExceeded) && err == nil {
				err = globalCtx.Err()
			}
			outcomes <- Outcome[T]{SourceID: task.SourceID, Value: value, Err: err, Duration: time.Since(start)}
		}
	}

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	sendDone := false
	for _, task := range tasks {
		if sendDone {
			break
		}
		select {
		case <-globalCtx.Done():
			sendDone = true
		case jobs <- task:
		}
	}
	close(jobs)
	wg.Wait()
	close(outcomes)

	results := make([]Outcome[T], 0, len(tasks))
	for outcome := range outcomes {
		results = append(results, outcome)
	}
	return results
}
