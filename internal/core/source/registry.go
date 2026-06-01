package source

import (
	"fmt"
	"sort"
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
	Modes            []Mode
	Priority         int
	Marketplace      bool
	AffiliateCapable bool
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

// CollectionIDs returns sorted collection IDs for discovery tools.
func (r Registry) CollectionIDs() []string {
	ids := make([]string, 0, len(r.Collections))
	for id := range r.Collections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
