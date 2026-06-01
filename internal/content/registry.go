package content

import "github.com/ralleh-ai/ralleh-mcp/internal/core/source"

// DefaultRegistry is the curated content/news/research source registry. It is
// for content lookup and synthesis, not arbitrary web crawling.
func DefaultRegistry() source.Registry {
	sources := map[string]source.Source{
		"reuters":             {ID: "reuters", Name: "Reuters", Domains: []string{"reuters.com"}, SearchTemplate: "https://www.reuters.com/site-search/?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 92},
		"ap_news":             {ID: "ap_news", Name: "Associated Press", Domains: []string{"apnews.com"}, SearchTemplate: "https://apnews.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 90},
		"bbc":                 {ID: "bbc", Name: "BBC", Domains: []string{"bbc.com"}, SearchTemplate: "https://www.bbc.co.uk/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 88},
		"npr":                 {ID: "npr", Name: "NPR", Domains: []string{"npr.org"}, SearchTemplate: "https://www.npr.org/search?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"guardian":            {ID: "guardian", Name: "The Guardian", Domains: []string{"theguardian.com"}, SearchTemplate: "https://www.theguardian.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"axios":               {ID: "axios", Name: "Axios", Domains: []string{"axios.com"}, SearchTemplate: "https://www.axios.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"cnbc":                {ID: "cnbc", Name: "CNBC", Domains: []string{"cnbc.com"}, SearchTemplate: "https://www.cnbc.com/search/?query={query}&qsearchterm={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"marketwatch":         {ID: "marketwatch", Name: "MarketWatch", Domains: []string{"marketwatch.com"}, SearchTemplate: "https://www.marketwatch.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80},
		"yahoo_finance":       {ID: "yahoo_finance", Name: "Yahoo Finance", Domains: []string{"finance.yahoo.com"}, SearchTemplate: "https://finance.yahoo.com/lookup?s={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 86},
		"sec_edgar":           {ID: "sec_edgar", Name: "SEC EDGAR", Domains: []string{"sec.gov"}, SearchTemplate: "https://www.sec.gov/edgar/search/#/q={query}", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 82},
		"espn":                {ID: "espn", Name: "ESPN", Domains: []string{"espn.com"}, SearchTemplate: "https://www.espn.com/search/_/q/{query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 86},
		"cbssports":           {ID: "cbssports", Name: "CBS Sports", Domains: []string{"cbssports.com"}, SearchTemplate: "https://www.cbssports.com/search/?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"sports_illustrated":  {ID: "sports_illustrated", Name: "Sports Illustrated", Domains: []string{"si.com"}, SearchTemplate: "https://www.si.com/search?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 74},
		"variety":             {ID: "variety", Name: "Variety", Domains: []string{"variety.com"}, SearchTemplate: "https://variety.com/results/#?q={query}", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 82},
		"hollywood_reporter":  {ID: "hollywood_reporter", Name: "The Hollywood Reporter", Domains: []string{"hollywoodreporter.com"}, SearchTemplate: "https://www.hollywoodreporter.com/?s={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80},
		"deadline":            {ID: "deadline", Name: "Deadline", Domains: []string{"deadline.com"}, SearchTemplate: "https://deadline.com/?s={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"sciencedaily":        {ID: "sciencedaily", Name: "ScienceDaily", Domains: []string{"sciencedaily.com"}, SearchTemplate: "https://www.sciencedaily.com/search/?keyword={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 86},
		"nature":              {ID: "nature", Name: "Nature", Domains: []string{"nature.com"}, SearchTemplate: "https://www.nature.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"scientific_american": {ID: "scientific_american", Name: "Scientific American", Domains: []string{"scientificamerican.com"}, SearchTemplate: "https://www.scientificamerican.com/search/?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80},
		"arxiv":               {ID: "arxiv", Name: "arXiv", Domains: []string{"arxiv.org"}, SearchTemplate: "https://arxiv.org/search/?query={query}&searchtype=all", Modes: []source.Mode{source.ModeHTMLFetch}, Priority: 92},
		"pubmed":              {ID: "pubmed", Name: "PubMed", Domains: []string{"pubmed.ncbi.nlm.nih.gov"}, SearchTemplate: "https://pubmed.ncbi.nlm.nih.gov/?term={query}", Modes: []source.Mode{source.ModeHTMLFetch}, Priority: 92},
		"semantic_scholar":    {ID: "semantic_scholar", Name: "Semantic Scholar", Domains: []string{"semanticscholar.org"}, SearchTemplate: "https://www.semanticscholar.org/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 88},
		"google_scholar":      {ID: "google_scholar", Name: "Google Scholar", Domains: []string{"scholar.google.com"}, SearchTemplate: "https://scholar.google.com/scholar?q={query}", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"hacker_news":         {ID: "hacker_news", Name: "Hacker News", Domains: []string{"hn.algolia.com"}, SearchTemplate: "https://hn.algolia.com/?q={query}", Modes: []source.Mode{source.ModeHTMLFetch}, Priority: 84},
		"the_verge":           {ID: "the_verge", Name: "The Verge", Domains: []string{"theverge.com"}, SearchTemplate: "https://www.theverge.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 80},
		"techcrunch":          {ID: "techcrunch", Name: "TechCrunch", Domains: []string{"techcrunch.com"}, SearchTemplate: "https://techcrunch.com/search/{query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
	}
	collections := map[string]source.Collection{
		"breaking_news":  {ID: "breaking_news", Label: "Breaking news", Description: "Current-event lookup across high-signal news sources.", DefaultSources: []string{"reuters", "ap_news", "bbc", "npr", "guardian"}, ExtendedSources: []string{"axios"}, MaxSources: 5},
		"stocks_markets": {ID: "stocks_markets", Label: "Stocks & markets", Description: "Market, company, filing, and finance lookup.", DefaultSources: []string{"yahoo_finance", "cnbc", "marketwatch", "reuters", "sec_edgar"}, MaxSources: 5},
		"sports":         {ID: "sports", Label: "Sports", Description: "Sports news, teams, players, leagues, and events.", DefaultSources: []string{"espn", "cbssports", "sports_illustrated", "ap_news"}, MaxSources: 4},
		"entertainment":  {ID: "entertainment", Label: "Entertainment", Description: "Film, TV, media, celebrities, streaming, box office, and industry news.", DefaultSources: []string{"variety", "hollywood_reporter", "deadline", "bbc"}, MaxSources: 4},
		"science":        {ID: "science", Label: "Science", Description: "Science reporting, discoveries, and public-facing scientific news.", DefaultSources: []string{"sciencedaily", "nature", "scientific_american", "bbc"}, ExtendedSources: []string{"reuters"}, MaxSources: 4},
		"research":       {ID: "research", Label: "Research", Description: "Academic papers, biomedical literature, preprints, and scholarly lookup.", DefaultSources: []string{"pubmed", "arxiv", "semantic_scholar"}, ExtendedSources: []string{"google_scholar", "nature"}, MaxSources: 4},
		"technology":     {ID: "technology", Label: "Technology", Description: "Technology news, startups, developer topics, AI, and internet culture.", DefaultSources: []string{"hacker_news", "the_verge", "techcrunch", "reuters"}, ExtendedSources: []string{"bbc"}, MaxSources: 4},
	}
	return source.Registry{Sources: sources, Collections: collections}
}
