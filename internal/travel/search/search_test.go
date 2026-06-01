package search

import (
	"context"
	"testing"
)

func TestSearchReturnsFakeFlightResults(t *testing.T) {
	resp, err := Search(context.Background(), Request{Origin: "MCO", Destination: "LAS", DepartDate: "2026-07-12", Collection: "us_domestic_flights", PreferredSources: []string{"duffel", "random_ota"}, BudgetProfile: "fast"}, FakeAdapter{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" || len(resp.Results) == 0 {
		t.Fatalf("expected flight options, got %+v", resp)
	}
	if len(resp.SourcePlan.RejectedSources) != 1 || resp.SourcePlan.RejectedSources[0] != "random_ota" {
		t.Fatalf("expected random OTA rejection, got %+v", resp.SourcePlan.RejectedSources)
	}
	if resp.Capabilities["canBook"] {
		t.Fatal("travel search must not expose booking capability")
	}
}
