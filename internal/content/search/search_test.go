package search

import (
	"context"
	"testing"
)

func TestContentSearchReturnsFakeResultsAndRejectsUnknownSource(t *testing.T) {
	resp, err := Search(context.Background(), Request{Query: "ai chips", Collection: "technology", PreferredSources: []string{"hacker_news", "random_blog"}, BudgetProfile: "fast"}, FakeAdapter{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" || len(resp.Results) == 0 {
		t.Fatalf("expected content results, got %+v", resp)
	}
	if len(resp.SourcePlan.RejectedSources) != 1 || resp.SourcePlan.RejectedSources[0] != "random_blog" {
		t.Fatalf("expected random source rejection, got %+v", resp.SourcePlan.RejectedSources)
	}
	if resp.Capabilities["canCrawlArbitraryWebsites"] {
		t.Fatal("content search must not expose arbitrary crawl capability")
	}
}
