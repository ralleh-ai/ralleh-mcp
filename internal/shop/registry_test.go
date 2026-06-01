package shop

import "testing"

func TestToolsCollectionUsesCuratedSourcesOnly(t *testing.T) {
	reg := DefaultRegistry()
	sources, rejected, err := reg.ResolveCollection("tools", []string{"harbor_freight", "random_site"}, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(rejected) != 1 || rejected[0] != "random_site" {
		t.Fatalf("expected random source rejection, got %#v", rejected)
	}
	if len(sources) == 0 || sources[0].ID != "harbor_freight" {
		t.Fatalf("preferred curated source should be first, got %#v", sources)
	}
	for _, s := range sources {
		if s.ID == "random_site" {
			t.Fatal("random source leaked into resolved plan")
		}
	}
}

func TestRegistrySearchTemplates(t *testing.T) {
	reg := DefaultRegistry()
	for id, src := range reg.Sources {
		if src.SearchTemplate == "" {
			t.Fatalf("source %s missing search template", id)
		}
		url, err := src.SearchURL("cordless drill 20v")
		if err != nil {
			t.Fatalf("source %s search URL: %v", id, err)
		}
		if url == src.SearchTemplate {
			t.Fatalf("source %s template was not expanded", id)
		}
	}
}

func TestExpandedCollectionsExist(t *testing.T) {
	reg := DefaultRegistry()
	for _, id := range []string{"tools", "office", "electronics", "clothing", "marketplaces", "auto", "toys", "gifts"} {
		if _, ok := reg.Collections[id]; !ok {
			t.Fatalf("missing collection %s", id)
		}
	}
}
