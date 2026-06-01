package executor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
)

func TestRunBoundedReturnsPartials(t *testing.T) {
	b := budget.Budget{GlobalTimeout: 200 * time.Millisecond, PerSourceTimeout: 50 * time.Millisecond, MaxSources: 3, MaxConcurrency: 2}
	tasks := []Task[string]{
		{SourceID: "fast", Run: func(context.Context) (string, error) { return "ok", nil }},
		{SourceID: "fail", Run: func(context.Context) (string, error) { return "", errors.New("boom") }},
		{SourceID: "slow", Run: func(ctx context.Context) (string, error) { <-ctx.Done(); return "", ctx.Err() }},
	}
	out := RunBounded(context.Background(), b, tasks)
	if len(out) != 3 {
		t.Fatalf("expected 3 outcomes, got %d", len(out))
	}
	seenFast := false
	seenErr := false
	for _, item := range out {
		if item.SourceID == "fast" && item.Value == "ok" {
			seenFast = true
		}
		if item.Err != nil {
			seenErr = true
		}
	}
	if !seenFast || !seenErr {
		t.Fatalf("expected partial success and errors, got %#v", out)
	}
}

func TestRunBoundedClampsMaxSources(t *testing.T) {
	b := budget.Budget{GlobalTimeout: time.Second, PerSourceTimeout: time.Second, MaxSources: 1, MaxConcurrency: 4}
	out := RunBounded(context.Background(), b, []Task[int]{
		{SourceID: "one", Run: func(context.Context) (int, error) { return 1, nil }},
		{SourceID: "two", Run: func(context.Context) (int, error) { return 2, nil }},
	})
	if len(out) != 1 {
		t.Fatalf("expected max source clamp to 1, got %d", len(out))
	}
}
