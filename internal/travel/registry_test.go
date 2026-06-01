package travel

import (
	"strings"
	"testing"
)

func TestTravelRejectsUnknownSourceIDs(t *testing.T) {
	reg := DefaultRegistry()
	sources, rejected, err := reg.ResolveCollection("us_domestic_flights", []string{"duffel", "random_ota"}, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) == 0 || sources[0].ID != "duffel" {
		t.Fatalf("expected preferred known source first, got %#v", sources)
	}
	if len(rejected) != 1 || rejected[0] != "random_ota" {
		t.Fatalf("expected unknown OTA rejection, got %#v", rejected)
	}
}

func TestTravelSearchTemplatesWhenConfigured(t *testing.T) {
	reg := DefaultRegistry()
	for id, src := range reg.Sources {
		if src.SearchTemplate == "" {
			continue
		}
		if !strings.Contains(src.SearchTemplate, "{query}") {
			continue
		}
		url, err := src.SearchURL("Orlando")
		if err != nil {
			t.Fatalf("source %s search URL: %v", id, err)
		}
		if url == src.SearchTemplate {
			t.Fatalf("source %s template was not expanded", id)
		}
	}
}

func TestHotelsCollectionExists(t *testing.T) {
	reg := DefaultRegistry()
	if _, ok := reg.Collections["hotels"]; !ok {
		t.Fatal("missing hotels collection")
	}
}
