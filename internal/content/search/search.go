package search

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/content"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/executor"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/result"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
)

type Request struct {
	Query            string         `json:"query"`
	Collection       string         `json:"collection"`
	PreferredSources []string       `json:"preferredSources,omitempty"`
	BudgetProfile    budget.Profile `json:"budgetProfile,omitempty"`
}

type Item struct {
	Title       string            `json:"title"`
	SourceID    string            `json:"sourceId"`
	SourceName  string            `json:"sourceName"`
	URL         string            `json:"url"`
	Summary     string            `json:"summary"`
	PublishedAt string            `json:"publishedAt,omitempty"`
	Confidence  float64           `json:"confidence"`
	Evidence    map[string]string `json:"evidence"`
}

type Response struct {
	Status       string            `json:"status"`
	Query        string            `json:"query"`
	Collection   string            `json:"collection"`
	SearchedAt   time.Time         `json:"searchedAt"`
	SourcePlan   result.SourcePlan `json:"sourcePlan"`
	Results      []Item            `json:"results"`
	Capabilities map[string]bool   `json:"capabilities"`
}

type Adapter interface {
	Search(context.Context, source.Source, Request, budget.Budget) ([]Item, error)
}

func Search(ctx context.Context, req Request, adapter Adapter) (Response, error) {
	if strings.TrimSpace(req.Query) == "" {
		return Response{}, fmt.Errorf("query is required")
	}
	if strings.TrimSpace(req.Collection) == "" {
		return Response{}, fmt.Errorf("collection is required")
	}
	if adapter == nil {
		adapter = FakeAdapter{}
	}
	b, err := budget.Resolve(req.BudgetProfile)
	if err != nil {
		return Response{}, err
	}
	reg := content.DefaultRegistry()
	sources, rejected, err := reg.ResolveCollection(req.Collection, req.PreferredSources, b.MaxSources)
	if err != nil {
		return Response{}, err
	}
	accepted := make([]string, 0, len(sources))
	tasks := make([]executor.Task[[]Item], 0, len(sources))
	for _, src := range sources {
		src := src
		accepted = append(accepted, src.ID)
		tasks = append(tasks, executor.Task[[]Item]{SourceID: src.ID, Run: func(ctx context.Context) ([]Item, error) { return adapter.Search(ctx, src, req, b) }})
	}
	outcomes := executor.RunBounded(ctx, b, tasks)
	items := []Item{}
	diagnostics := []result.SourceDiagnostic{}
	for _, outcome := range outcomes {
		status, errType, errText := "success", "", ""
		if outcome.Err != nil {
			status, errType, errText = "error", "source_error", outcome.Err.Error()
		}
		diagnostics = append(diagnostics, result.SourceDiagnostic{SourceID: outcome.SourceID, Status: status, Mode: "fake_upstream", Duration: outcome.Duration, ErrorType: errType, Error: errText, ResultCount: len(outcome.Value)})
		items = append(items, outcome.Value...)
	}
	status := "ok"
	if len(items) == 0 {
		status = "empty"
	}
	return Response{Status: status, Query: req.Query, Collection: req.Collection, SearchedAt: time.Now().UTC(), SourcePlan: result.SourcePlan{Collection: req.Collection, RequestedSources: req.PreferredSources, AcceptedSources: accepted, RejectedSources: rejected, BudgetProfile: string(b.Profile), Diagnostics: diagnostics}, Results: items, Capabilities: map[string]bool{"canSearchContent": true, "canSummarize": true, "canCrawlArbitraryWebsites": false}}, nil
}

type FakeAdapter struct{}

func (FakeAdapter) Search(ctx context.Context, src source.Source, req Request, b budget.Budget) ([]Item, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	count := b.MaxResultsPerSource
	if count > 2 {
		count = 2
	}
	if count <= 0 {
		count = 1
	}
	out := make([]Item, 0, count)
	searchURL, err := src.SearchURL(req.Query)
	if err != nil {
		return nil, err
	}
	for i := 1; i <= count; i++ {
		id := stableID(src.ID, req.Query, i)
		out = append(out, Item{Title: fmt.Sprintf("%s item %d for %s", src.Name, i, req.Query), SourceID: src.ID, SourceName: src.Name, URL: fmt.Sprintf("%s#ralleh-smoke-%s", searchURL, id), Summary: fmt.Sprintf("Deterministic smoke result from %s for %s. Use real adapters for live content extraction.", src.Name, req.Query), Confidence: 0.89, Evidence: map[string]string{"mode": "fake_upstream", "searchUrl": searchURL}})
	}
	return out, nil
}

func stableID(sourceID, query string, n int) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%s:%s:%d", sourceID, strings.ToLower(query), n)))
	return hex.EncodeToString(h[:])[:12]
}
