package search

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/executor"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/result"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
	"github.com/ralleh-ai/ralleh-mcp/internal/travel"
)

// Request is the v1 travel research contract. It is research-only: no booking,
// no payment, no passenger PII.
type Request struct {
	Origin           string         `json:"origin"`
	Destination      string         `json:"destination"`
	DepartDate       string         `json:"departDate"`
	ReturnDate       string         `json:"returnDate,omitempty"`
	Collection       string         `json:"collection"`
	PreferredSources []string       `json:"preferredSources,omitempty"`
	BudgetProfile    budget.Profile `json:"budgetProfile,omitempty"`
}

type FlightOption struct {
	ID                   string   `json:"id"`
	Provider             string   `json:"provider"`
	SourceID             string   `json:"sourceId"`
	Price                float64  `json:"price"`
	Currency             string   `json:"currency"`
	Origin               string   `json:"origin"`
	Destination          string   `json:"destination"`
	DepartDate           string   `json:"departDate"`
	ReturnDate           string   `json:"returnDate,omitempty"`
	Stops                int      `json:"stops"`
	TotalDurationMinutes int      `json:"totalDurationMinutes"`
	BookingURL           string   `json:"bookingUrl"`
	Confidence           float64  `json:"confidence"`
	Warnings             []string `json:"warnings,omitempty"`
}

type Response struct {
	Status       string            `json:"status"`
	SearchedAt   time.Time         `json:"searchedAt"`
	SourcePlan   result.SourcePlan `json:"sourcePlan"`
	Results      []FlightOption    `json:"results"`
	Capabilities map[string]bool   `json:"capabilities"`
}

type Adapter interface {
	Search(context.Context, source.Source, Request, budget.Budget) ([]FlightOption, error)
}

func Search(ctx context.Context, req Request, adapter Adapter) (Response, error) {
	if strings.TrimSpace(req.Origin) == "" || strings.TrimSpace(req.Destination) == "" || strings.TrimSpace(req.DepartDate) == "" {
		return Response{}, fmt.Errorf("origin, destination, and departDate are required")
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
	reg := travel.DefaultRegistry()
	sources, rejected, err := reg.ResolveCollection(req.Collection, req.PreferredSources, b.MaxSources)
	if err != nil {
		return Response{}, err
	}
	accepted := make([]string, 0, len(sources))
	tasks := make([]executor.Task[[]FlightOption], 0, len(sources))
	for _, src := range sources {
		src := src
		accepted = append(accepted, src.ID)
		tasks = append(tasks, executor.Task[[]FlightOption]{SourceID: src.ID, Run: func(ctx context.Context) ([]FlightOption, error) { return adapter.Search(ctx, src, req, b) }})
	}
	outcomes := executor.RunBounded(ctx, b, tasks)
	options := []FlightOption{}
	diagnostics := []result.SourceDiagnostic{}
	for _, outcome := range outcomes {
		status := "success"
		errType := ""
		errText := ""
		if outcome.Err != nil {
			status = "error"
			errType = "source_error"
			errText = outcome.Err.Error()
		}
		diagnostics = append(diagnostics, result.SourceDiagnostic{SourceID: outcome.SourceID, Status: status, Mode: "fake_upstream", Duration: outcome.Duration, ErrorType: errType, Error: errText, ResultCount: len(outcome.Value)})
		options = append(options, outcome.Value...)
	}
	status := "ok"
	if len(options) == 0 {
		status = "empty"
	}
	return Response{Status: status, SearchedAt: time.Now().UTC(), SourcePlan: result.SourcePlan{Collection: req.Collection, RequestedSources: req.PreferredSources, AcceptedSources: accepted, RejectedSources: rejected, BudgetProfile: string(b.Profile), Diagnostics: diagnostics}, Results: options, Capabilities: map[string]bool{"canSearchFlights": true, "canBook": false, "canUseCreditCard": false, "canEnterPassengerInfo": false}}, nil
}

type FakeAdapter struct{}

func (FakeAdapter) Search(ctx context.Context, src source.Source, req Request, b budget.Budget) ([]FlightOption, error) {
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
	out := make([]FlightOption, 0, count)
	for i := 1; i <= count; i++ {
		id := stableID(src.ID, req.Origin, req.Destination, req.DepartDate, i)
		warnings := []string{"research_only_no_booking"}
		if src.ID == "kiwi" || src.ID == "skyscanner" {
			warnings = append(warnings, "verify baggage and fare rules before booking manually")
		}
		out = append(out, FlightOption{ID: id, Provider: src.Name, SourceID: src.ID, Price: float64(129 + i*25 + len(src.ID)), Currency: "USD", Origin: strings.ToUpper(req.Origin), Destination: strings.ToUpper(req.Destination), DepartDate: req.DepartDate, ReturnDate: req.ReturnDate, Stops: i - 1, TotalDurationMinutes: 140 + i*55 + len(src.ID), BookingURL: fmt.Sprintf("https://www.%s/ralleh-smoke/flights/%s", src.Domains[0], id), Confidence: 0.88, Warnings: warnings})
	}
	return out, nil
}

func stableID(parts ...any) string {
	h := sha1.Sum([]byte(fmt.Sprint(parts...)))
	return hex.EncodeToString(h[:])[:12]
}
