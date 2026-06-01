package shop

import "github.com/ralleh-ai/ralleh-mcp/internal/core/source"

// DefaultRegistry is the first curated shopping source registry. Raw website
// URLs are not accepted by v1 search APIs; callers choose collections and known
// source IDs only.
func DefaultRegistry() source.Registry {
	sources := map[string]source.Source{
		"home_depot":     {ID: "home_depot", Name: "Home Depot", Domains: []string{"homedepot.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 95},
		"lowes":          {ID: "lowes", Name: "Lowe's", Domains: []string{"lowes.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 92},
		"menards":        {ID: "menards", Name: "Menards", Domains: []string{"menards.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 75},
		"ace_hardware":   {ID: "ace_hardware", Name: "Ace Hardware", Domains: []string{"acehardware.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80},
		"tractor_supply": {ID: "tractor_supply", Name: "Tractor Supply", Domains: []string{"tractorsupply.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"harbor_freight": {ID: "harbor_freight", Name: "Harbor Freight", Domains: []string{"harborfreight.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 88},
		"best_buy":       {ID: "best_buy", Name: "Best Buy", Domains: []string{"bestbuy.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 92},
		"costco":         {ID: "costco", Name: "Costco", Domains: []string{"costco.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 82},
		"target":         {ID: "target", Name: "Target", Domains: []string{"target.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 88},
		"walmart":        {ID: "walmart", Name: "Walmart", Domains: []string{"walmart.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 86, Marketplace: true},
		"ebay":           {ID: "ebay", Name: "eBay", Domains: []string{"ebay.com"}, Modes: []source.Mode{source.ModeAPI, source.ModeHTMLFetch}, Priority: 80, Marketplace: true, AffiliateCapable: true},
		"etsy":           {ID: "etsy", Name: "Etsy", Domains: []string{"etsy.com"}, Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80, Marketplace: true},
		"amazon_api":     {ID: "amazon_api", Name: "Amazon Product Advertising API", Domains: []string{"amazon.com"}, Modes: []source.Mode{source.ModeAPI}, Priority: 70, Marketplace: true, AffiliateCapable: true},
	}
	collections := map[string]source.Collection{
		"tools": {
			ID: "tools", Label: "Tools & hardware", Description: "Tools, hardware, jobsite equipment, power tools, lawn/farm utility, and shop gear.",
			DefaultSources:  []string{"home_depot", "lowes", "harbor_freight", "ace_hardware", "tractor_supply", "ebay"},
			ExtendedSources: []string{"menards", "walmart", "amazon_api"}, MaxSources: 5,
		},
		"toys": {
			ID: "toys", Label: "Toys & games", Description: "Toys, games, kids gifts, and family entertainment.",
			DefaultSources:  []string{"target", "walmart", "best_buy", "costco", "ebay"},
			ExtendedSources: []string{"amazon_api"}, MaxSources: 5,
		},
		"gifts": {
			ID: "gifts", Label: "Gifts", Description: "General gifts, handmade goods, marketplaces, and broad retail options.",
			DefaultSources:  []string{"etsy", "target", "ebay", "walmart"},
			ExtendedSources: []string{"costco", "amazon_api"}, MaxSources: 5,
		},
		"electronics": {
			ID: "electronics", Label: "Electronics", Description: "Consumer electronics, computers, accessories, gaming hardware, and appliances-adjacent tech.",
			DefaultSources:  []string{"best_buy", "target", "walmart", "costco", "ebay"},
			ExtendedSources: []string{"amazon_api"}, MaxSources: 5,
		},
	}
	return source.Registry{Sources: sources, Collections: collections}
}
