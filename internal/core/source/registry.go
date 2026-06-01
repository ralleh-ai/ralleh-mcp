package source

import (
	"fmt"
	"math"
	"net/url"
	"sort"
	"strings"
)

// Mode describes how a known source may be queried.
type Mode string

const (
	ModeAPI           Mode = "api"
	ModeFeed          Mode = "feed"
	ModeHTMLFetch     Mode = "html_fetch"
	ModeBrowserVerify Mode = "browser_verify"
)

// Source describes one approved source adapter. IDs are stable API values used
// by LLM clients; raw URLs are intentionally not accepted in search requests.
type Source struct {
	ID               string
	Name             string
	Domains          []string
	SearchTemplate   string
	Modes            []Mode
	Priority         int
	Marketplace      bool
	AffiliateCapable bool
}

// Ranking explains why one curated source should be tried before another.
type Ranking struct {
	SourceID      string   `json:"sourceId"`
	Name          string   `json:"name"`
	Speed         int      `json:"speed"`
	Accuracy      int      `json:"accuracy"`
	Usefulness    int      `json:"usefulness"`
	Reliability   int      `json:"reliability"`
	ChallengeRisk int      `json:"challengeRisk"`
	Overall       float64  `json:"overall"`
	Modes         []Mode   `json:"modes"`
	Rationale     []string `json:"rationale"`
}

// Ranking returns a deterministic v0 quality score. It is a starting heuristic;
// live probes and adapter success metrics should feed these signals later.
func (s Source) Ranking() Ranking {
	speed := clampScore(s.Priority)
	accuracy := clampScore(s.Priority)
	usefulness := clampScore(s.Priority)
	reliability := clampScore(s.Priority)
	challenge := 30
	rationale := []string{}

	if hasMode(s.Modes, ModeAPI) {
		speed += 10
		accuracy += 10
		reliability += 10
		challenge -= 15
		rationale = append(rationale, "api-capable source")
	}
	if hasMode(s.Modes, ModeFeed) {
		speed += 8
		accuracy += 8
		reliability += 8
		challenge -= 10
		rationale = append(rationale, "feed-capable source")
	}
	if hasMode(s.Modes, ModeBrowserVerify) && !hasMode(s.Modes, ModeHTMLFetch) && !hasMode(s.Modes, ModeAPI) {
		speed -= 25
		reliability -= 15
		challenge += 25
		rationale = append(rationale, "browser verification required")
	}
	if s.SearchTemplate == "" {
		speed -= 15
		reliability -= 10
		challenge += 15
		rationale = append(rationale, "no direct search template")
	} else if strings.Contains(s.SearchTemplate, "{query}") {
		speed += 5
		reliability += 5
		rationale = append(rationale, "direct query URL template")
	}
	if s.Marketplace {
		usefulness += 6
		accuracy -= 4
		rationale = append(rationale, "marketplace breadth; verify sellers/condition")
	}
	if s.AffiliateCapable {
		usefulness += 3
		rationale = append(rationale, "affiliate capable")
	}

	speed = clampScore(speed)
	accuracy = clampScore(accuracy)
	usefulness = clampScore(usefulness)
	reliability = clampScore(reliability)
	challenge = clampScore(challenge)
	overall := 0.25*float64(speed) + 0.35*float64(accuracy) + 0.30*float64(usefulness) + 0.10*float64(reliability) - 0.10*float64(challenge)
	overall = math.Round(overall*10) / 10
	return Ranking{SourceID: s.ID, Name: s.Name, Speed: speed, Accuracy: accuracy, Usefulness: usefulness, Reliability: reliability, ChallengeRisk: challenge, Overall: overall, Modes: s.Modes, Rationale: rationale}
}

// SearchURL builds a source search URL from the source-owned template. Templates
// must contain {query}. The query is path-escaped for path placeholders and
// query-escaped for query-string placeholders.
func (s Source) SearchURL(query string) (string, error) {
	if strings.TrimSpace(s.SearchTemplate) == "" {
		return "", fmt.Errorf("source %q has no search template", s.ID)
	}
	if !strings.Contains(s.SearchTemplate, "{query}") {
		return "", fmt.Errorf("source %q search template missing {query}", s.ID)
	}
	encoded := url.QueryEscape(query)
	if strings.Contains(s.SearchTemplate, "/{query}") || strings.Contains(s.SearchTemplate, "={query}") == false && strings.HasSuffix(s.SearchTemplate, "{query}") {
		encoded = url.PathEscape(query)
	}
	return strings.ReplaceAll(s.SearchTemplate, "{query}", encoded), nil
}

// Collection maps a product/travel category to approved source IDs.
type Collection struct {
	ID              string
	Label           string
	Description     string
	DefaultSources  []string
	ExtendedSources []string
	MaxSources      int
}

// Registry stores known sources and curated collections.
type Registry struct {
	Sources     map[string]Source
	Collections map[string]Collection
}

// ResolveCollection validates a collection and optional preferred source IDs,
// then returns the ordered approved source list. Unknown IDs are rejected.
func (r Registry) ResolveCollection(collectionID string, preferred []string, maxSources int) ([]Source, []string, error) {
	collection, ok := r.Collections[collectionID]
	if !ok {
		return nil, nil, fmt.Errorf("unknown collection %q", collectionID)
	}
	allowed := map[string]bool{}
	ordered := append([]string{}, collection.DefaultSources...)
	ordered = append(ordered, collection.ExtendedSources...)
	for _, id := range ordered {
		allowed[id] = true
	}

	seen := map[string]bool{}
	selectedIDs := make([]string, 0, maxSources)
	rejected := []string{}
	appendID := func(id string) {
		if len(selectedIDs) >= maxSources || seen[id] {
			return
		}
		if _, exists := r.Sources[id]; !exists || !allowed[id] {
			rejected = append(rejected, id)
			return
		}
		seen[id] = true
		selectedIDs = append(selectedIDs, id)
	}

	for _, id := range preferred {
		appendID(id)
	}
	for _, id := range collection.DefaultSources {
		appendID(id)
	}

	selected := make([]Source, 0, len(selectedIDs))
	for _, id := range selectedIDs {
		selected = append(selected, r.Sources[id])
	}
	return selected, rejected, nil
}

// RankCollection returns all default and extended sources for a collection sorted
// by deterministic quality score, highest first.
func (r Registry) RankCollection(collectionID string) ([]Ranking, error) {
	collection, ok := r.Collections[collectionID]
	if !ok {
		return nil, fmt.Errorf("unknown collection %q", collectionID)
	}
	ids := append([]string{}, collection.DefaultSources...)
	ids = append(ids, collection.ExtendedSources...)
	seen := map[string]bool{}
	rankings := make([]Ranking, 0, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		src, ok := r.Sources[id]
		if !ok {
			continue
		}
		rankings = append(rankings, src.Ranking())
	}
	sort.SliceStable(rankings, func(i, j int) bool {
		if rankings[i].Overall == rankings[j].Overall {
			return rankings[i].SourceID < rankings[j].SourceID
		}
		return rankings[i].Overall > rankings[j].Overall
	})
	return rankings, nil
}

// RankAll returns rankings for every collection.
func (r Registry) RankAll() map[string][]Ranking {
	out := map[string][]Ranking{}
	for _, id := range r.CollectionIDs() {
		rankings, err := r.RankCollection(id)
		if err == nil {
			out[id] = rankings
		}
	}
	return out
}

func hasMode(modes []Mode, mode Mode) bool {
	for _, item := range modes {
		if item == mode {
			return true
		}
	}
	return false
}

func clampScore(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// CollectionIDs returns sorted collection IDs for discovery tools.
func (r Registry) CollectionIDs() []string {
	ids := make([]string, 0, len(r.Collections))
	for id := range r.Collections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
