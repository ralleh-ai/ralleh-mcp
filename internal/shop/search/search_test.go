package search

import (
	"context"
	"testing"
)

func TestSearchRejectsUnknownSourcesAndReturnsFakeResults(t *testing.T) {
	resp, err := Search(context.Background(), Request{Query: "cordless drill", Collection: "tools", PreferredSources: []string{"harbor_freight", "random_site"}, BudgetProfile: "fast"}, FakeAdapter{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" || len(resp.Results) == 0 {
		t.Fatalf("expected ok results, got %+v", resp)
	}
	if len(resp.SourcePlan.RejectedSources) != 1 || resp.SourcePlan.RejectedSources[0] != "random_site" {
		t.Fatalf("expected random source rejection, got %+v", resp.SourcePlan.RejectedSources)
	}
	if resp.SourcePlan.AcceptedSources[0] != "harbor_freight" {
		t.Fatalf("expected preferred source first, got %+v", resp.SourcePlan.AcceptedSources)
	}
}

func TestSearchRequiresQueryAndCollection(t *testing.T) {
	if _, err := Search(context.Background(), Request{Collection: "tools"}, nil); err == nil {
		t.Fatal("expected missing query error")
	}
	if _, err := Search(context.Background(), Request{Query: "drill"}, nil); err == nil {
		t.Fatal("expected missing collection error")
	}
}
