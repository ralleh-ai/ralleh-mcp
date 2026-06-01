package travel

import "github.com/ralleh-ai/ralleh-mcp/internal/core/source"

// DefaultRegistry is the curated v1 travel source registry. Travel is
// research-only: no booking, no card use, no passenger PII.
func DefaultRegistry() source.Registry {
	sources := map[string]source.Source{
		"amadeus":         {ID: "amadeus", Name: "Amadeus", Domains: []string{"amadeus.com"}, Modes: []source.Mode{source.ModeAPI}, Priority: 98},
		"duffel":          {ID: "duffel", Name: "Duffel", Domains: []string{"duffel.com"}, Modes: []source.Mode{source.ModeAPI}, Priority: 96},
		"kiwi":            {ID: "kiwi", Name: "Kiwi", Domains: []string{"kiwi.com"}, Modes: []source.Mode{source.ModeAPI, source.ModeHTMLFetch}, Priority: 88},
		"skyscanner":      {ID: "skyscanner", Name: "Skyscanner", Domains: []string{"skyscanner.com"}, Modes: []source.Mode{source.ModeAPI, source.ModeBrowserVerify}, Priority: 90},
		"delta":           {ID: "delta", Name: "Delta", Domains: []string{"delta.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"united":          {ID: "united", Name: "United", Domains: []string{"united.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"american":        {ID: "american", Name: "American Airlines", Domains: []string{"aa.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"southwest":       {ID: "southwest", Name: "Southwest", Domains: []string{"southwest.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"jetblue":         {ID: "jetblue", Name: "JetBlue", Domains: []string{"jetblue.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"alaska":          {ID: "alaska", Name: "Alaska Airlines", Domains: []string{"alaskaair.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"frontier":        {ID: "frontier", Name: "Frontier", Domains: []string{"flyfrontier.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 55},
		"spirit":          {ID: "spirit", Name: "Spirit", Domains: []string{"spirit.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 55},
		"british_airways": {ID: "british_airways", Name: "British Airways", Domains: []string{"britishairways.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"lufthansa":       {ID: "lufthansa", Name: "Lufthansa", Domains: []string{"lufthansa.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"air_france":      {ID: "air_france", Name: "Air France", Domains: []string{"airfrance.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"klm":             {ID: "klm", Name: "KLM", Domains: []string{"klm.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
	}
	collections := map[string]source.Collection{
		"us_domestic_flights": {
			ID: "us_domestic_flights", Label: "US domestic flights", Description: "US domestic flight research via API providers first, airline direct as limited verification fallback.",
			DefaultSources:  []string{"amadeus", "duffel", "kiwi", "skyscanner"},
			ExtendedSources: []string{"delta", "united", "american", "southwest", "jetblue", "alaska"}, MaxSources: 4,
		},
		"international_flights": {
			ID: "international_flights", Label: "International flights", Description: "International flight research via API/aggregator sources first.",
			DefaultSources:  []string{"amadeus", "duffel", "skyscanner", "kiwi"},
			ExtendedSources: []string{"united", "delta", "american", "british_airways", "lufthansa", "air_france", "klm"}, MaxSources: 4,
		},
		"budget_flights": {
			ID: "budget_flights", Label: "Budget flights", Description: "Budget and low-cost carrier focused flight research with strong fare-trap warnings.",
			DefaultSources:  []string{"kiwi", "skyscanner"},
			ExtendedSources: []string{"frontier", "spirit", "southwest"}, MaxSources: 4,
		},
	}
	return source.Registry{Sources: sources, Collections: collections}
}
