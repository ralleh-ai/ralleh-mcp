package content

import "testing"

func TestContentCollectionsExist(t *testing.T) {
	reg := DefaultRegistry()
	for _, id := range []string{"web", "breaking_news", "news_brief", "rss", "stocks_markets", "sports", "weather", "government", "community", "reviews_consensus", "entertainment", "science", "research", "technology"} {
		if _, ok := reg.Collections[id]; !ok {
			t.Fatalf("missing collection %s", id)
		}
	}
}

func TestContentSourceTemplates(t *testing.T) {
	reg := DefaultRegistry()
	for id, src := range reg.Sources {
		if src.SearchTemplate == "" {
			t.Fatalf("content source %s missing search template", id)
		}
		if _, err := src.SearchURL("ai chips"); err != nil {
			t.Fatalf("content source %s search URL: %v", id, err)
		}
	}
}
