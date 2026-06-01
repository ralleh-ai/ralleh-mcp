package source

import "testing"

func TestSearchURLQueryEncoding(t *testing.T) {
	s := Source{ID: "walmart", SearchTemplate: "https://www.walmart.com/search?q={query}"}
	u, err := s.SearchURL("cordless drill 20v")
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://www.walmart.com/search?q=cordless+drill+20v" {
		t.Fatalf("unexpected URL: %s", u)
	}
}

func TestSearchURLPathEncoding(t *testing.T) {
	s := Source{ID: "home_depot", SearchTemplate: "https://www.homedepot.com/s/{query}"}
	u, err := s.SearchURL("cordless drill 20v")
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://www.homedepot.com/s/cordless%20drill%2020v" {
		t.Fatalf("unexpected URL: %s", u)
	}
}

func TestSearchURLRequiresTemplate(t *testing.T) {
	if _, err := (Source{ID: "x"}).SearchURL("drill"); err == nil {
		t.Fatal("expected missing template error")
	}
}

func TestRankCollectionOrdersByOverall(t *testing.T) {
	reg := Registry{
		Sources: map[string]Source{
			"slow": {ID: "slow", Name: "Slow", Priority: 30, Modes: []Mode{ModeBrowserVerify}},
			"fast": {ID: "fast", Name: "Fast", Priority: 80, SearchTemplate: "https://example.com/search?q={query}", Modes: []Mode{ModeAPI}},
		},
		Collections: map[string]Collection{"c": {ID: "c", DefaultSources: []string{"slow", "fast"}}},
	}
	ranks, err := reg.RankCollection("c")
	if err != nil {
		t.Fatal(err)
	}
	if len(ranks) != 2 || ranks[0].SourceID != "fast" {
		t.Fatalf("expected fast source first, got %#v", ranks)
	}
}
