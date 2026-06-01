package search

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/affiliate"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/executor"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/result"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
	"github.com/ralleh-ai/ralleh-mcp/internal/shop"
)

// Request is the v1 shop search contract. It accepts a curated collection and
// known source IDs only; it does not accept arbitrary websites.
type Request struct {
	Query            string         `json:"query"`
	Collection       string         `json:"collection"`
	PreferredSources []string       `json:"preferredSources,omitempty"`
	BudgetProfile    budget.Profile `json:"budgetProfile,omitempty"`
	AffiliateProfile string         `json:"affiliateProfile,omitempty"`
}

// Product is a normalized, LLM-ready product candidate.
type Product struct {
	Title        string            `json:"title"`
	Merchant     string            `json:"merchant"`
	SourceID     string            `json:"sourceId"`
	Price        float64           `json:"price"`
	Currency     string            `json:"currency"`
	Condition    string            `json:"condition"`
	Availability string            `json:"availability"`
	CanonicalURL string            `json:"canonicalUrl"`
	PresentedURL string            `json:"presentedUrl"`
	Affiliate    affiliate.Result  `json:"affiliate"`
	Confidence   float64           `json:"confidence"`
	Warnings     []string          `json:"warnings,omitempty"`
	Evidence     map[string]string `json:"evidence"`
}

// Response is the normalized search output.
type Response struct {
	Status              string                     `json:"status"`
	Query               string                     `json:"query"`
	Collection          string                     `json:"collection"`
	SearchedAt          time.Time                  `json:"searchedAt"`
	SourcePlan          result.SourcePlan          `json:"sourcePlan"`
	Results             []Product                  `json:"results"`
	AffiliateDisclosure result.AffiliateDisclosure `json:"affiliateDisclosure"`
	Capabilities        map[string]bool            `json:"capabilities"`
}

// Adapter searches one curated source.
type Adapter interface {
	Search(context.Context, source.Source, Request, budget.Budget) ([]Product, error)
}

// Search executes a bounded curated search using the provided adapter. In the
// current scaffold, FakeAdapter is used for deterministic smoke tests; real
// source adapters will plug into the same contract.
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
	reg := shop.DefaultRegistry()
	sources, rejected, err := reg.ResolveCollection(req.Collection, req.PreferredSources, b.MaxSources)
	if err != nil {
		return Response{}, err
	}
	accepted := make([]string, 0, len(sources))
	for _, src := range sources {
		accepted = append(accepted, src.ID)
	}
	tasks := make([]executor.Task[[]Product], 0, len(sources))
	for _, src := range sources {
		src := src
		tasks = append(tasks, executor.Task[[]Product]{SourceID: src.ID, Run: func(ctx context.Context) ([]Product, error) {
			return adapter.Search(ctx, src, req, b)
		}})
	}
	outcomes := executor.RunBounded(ctx, b, tasks)
	products := []Product{}
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
		products = append(products, outcome.Value...)
	}
	status := "ok"
	if len(products) == 0 {
		status = "empty"
	}
	return Response{
		Status:              status,
		Query:               req.Query,
		Collection:          req.Collection,
		SearchedAt:          time.Now().UTC(),
		SourcePlan:          result.SourcePlan{Collection: req.Collection, RequestedSources: req.PreferredSources, AcceptedSources: accepted, RejectedSources: rejected, BudgetProfile: string(b.Profile), Diagnostics: diagnostics},
		Results:             products,
		AffiliateDisclosure: result.AffiliateDisclosure{Required: hasAffiliate(products), Text: "Some links may be affiliate links. Ralleh may earn a commission at no extra cost to you."},
		Capabilities:        map[string]bool{"canSearch": true, "canCompare": false, "canVerify": false, "canPurchase": false, "canUseCreditCard": false},
	}, nil
}

func hasAffiliate(products []Product) bool {
	for _, p := range products {
		if p.Affiliate.Applied {
			return true
		}
	}
	return false
}

// FakeAdapter provides deterministic search output for tests and smoke checks.
type FakeAdapter struct{}

func (FakeAdapter) Search(ctx context.Context, src source.Source, req Request, b budget.Budget) ([]Product, error) {
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
	items := make([]Product, 0, count)
	for i := 1; i <= count; i++ {
		id := stableID(src.ID, req.Query, i)
		canonical := fmt.Sprintf("https://www.%s/ralleh-smoke/%s", src.Domains[0], id)
		aff := affiliate.Result{CanonicalURL: canonical, PresentedURL: canonical, Reason: "no_affiliate_rule_configured"}
		if src.ID == "ebay" {
			var err error
			aff, err = affiliate.ApplyQueryParam(canonical, affiliate.Rule{SourceID: "ebay", AllowedDomains: []string{"ebay.com"}, Param: "campid", Value: "ralleh-smoke", Enabled: true})
			if err != nil {
				return nil, err
			}
		}
		items = append(items, Product{Title: fmt.Sprintf("%s result %d for %s", src.Name, i, req.Query), Merchant: src.Name, SourceID: src.ID, Price: float64(49 + i*10 + len(src.ID)), Currency: "USD", Condition: "new", Availability: "in_stock", CanonicalURL: canonical, PresentedURL: aff.PresentedURL, Affiliate: aff, Confidence: 0.91, Evidence: map[string]string{"mode": "fake_upstream", "sourceDomain": src.Domains[0]}})
	}
	return items, nil
}

func stableID(sourceID, query string, n int) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%s:%s:%d", sourceID, strings.ToLower(query), n)))
	return hex.EncodeToString(h[:])[:12]
}
