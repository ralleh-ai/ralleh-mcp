package content

// ToolMap maps conceptual MCP tool names to curated content collections. These
// are not separate crawl permissions; they are safe aliases for known source sets.
func ToolMap() map[string]string {
	return map[string]string{
		"SearchWeb":          "web",
		"SearchNews":         "breaking_news",
		"GetBreakingNews":    "breaking_news",
		"GetNewsBrief":       "news_brief",
		"SearchRSS":          "rss",
		"SearchFinance":      "stocks_markets",
		"SearchStocks":       "stocks_markets",
		"SearchSports":       "sports",
		"SearchWeather":      "weather",
		"SearchGovernment":   "government",
		"SearchCommunity":    "community",
		"GetConsensus":       "reviews_consensus",
		"SummarizeConsensus": "reviews_consensus",
		"DetectConflicts":    "news_brief",
	}
}
