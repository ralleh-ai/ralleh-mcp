package travel

import "github.com/ralleh-ai/ralleh-mcp/internal/core/source"

// DefaultRegistry is the curated v1 travel source registry. Travel is
// research-only: no booking, no card use, no passenger PII.
func DefaultRegistry() source.Registry {
	sources := map[string]source.Source{
		"booking_com":        {ID: "booking_com", Name: "Booking.com", Domains: []string{"booking.com"}, SearchTemplate: "https://www.booking.com/searchresults.html?ss={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 76},
		"expedia_hotels":     {ID: "expedia_hotels", Name: "Expedia Hotels", Domains: []string{"expedia.com"}, SearchTemplate: "https://www.expedia.com/Hotel-Search?destination={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 76},
		"hotels_com":         {ID: "hotels_com", Name: "Hotels.com", Domains: []string{"hotels.com"}, SearchTemplate: "https://www.hotels.com/Hotel-Search?destination={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 74},
		"vrbo":               {ID: "vrbo", Name: "Vrbo", Domains: []string{"vrbo.com"}, SearchTemplate: "https://www.vrbo.com/search?destination={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 72},
		"kayak_flights":      {ID: "kayak_flights", Name: "Kayak Flights", Domains: []string{"kayak.com"}, SearchTemplate: "https://www.kayak.com/flights/{query}", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 68},
		"skyscanner":         {ID: "skyscanner", Name: "Skyscanner", Domains: []string{"skyscanner.com"}, SearchTemplate: "https://www.skyscanner.com/transport/flights/{query}", Modes: []source.Mode{source.ModeAPI, source.ModeBrowserVerify}, Priority: 90},
		"priceline_hotels":   {ID: "priceline_hotels", Name: "Priceline Hotels", Domains: []string{"priceline.com"}, SearchTemplate: "https://www.priceline.com/stay/#/search/?location={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 70},
		"travelocity_hotels": {ID: "travelocity_hotels", Name: "Travelocity Hotels", Domains: []string{"travelocity.com"}, SearchTemplate: "https://www.travelocity.com/Hotel-Search?destination={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 70},
		"google_flights":     {ID: "google_flights", Name: "Google Flights", Domains: []string{"google.com"}, SearchTemplate: "https://www.google.com/travel/flights", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 55},
		"momondo_flights":    {ID: "momondo_flights", Name: "Momondo Flights", Domains: []string{"momondo.com"}, SearchTemplate: "https://www.momondo.com/flight-search", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 62},
		"cheapoair_flights":  {ID: "cheapoair_flights", Name: "CheapOair Flights", Domains: []string{"cheapoair.com"}, SearchTemplate: "https://www.cheapoair.com/flights/search", Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 58},
		"recreation_gov":     {ID: "recreation_gov", Name: "Recreation.gov", Domains: []string{"recreation.gov"}, SearchTemplate: "https://www.recreation.gov/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"hipcamp":            {ID: "hipcamp", Name: "Hipcamp", Domains: []string{"hipcamp.com"}, SearchTemplate: "https://www.hipcamp.com/en-US/search?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 82},
		"the_dyrt":           {ID: "the_dyrt", Name: "The Dyrt", Domains: []string{"thedyrt.com"}, SearchTemplate: "https://thedyrt.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"campendium":         {ID: "campendium", Name: "Campendium", Domains: []string{"campendium.com"}, SearchTemplate: "https://www.campendium.com/search?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 76},
		"cruise_critic":      {ID: "cruise_critic", Name: "Cruise Critic", Domains: []string{"cruisecritic.com"}, SearchTemplate: "https://www.cruisecritic.com/search/?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 76},
		"vacations_to_go":    {ID: "vacations_to_go", Name: "Vacations To Go", Domains: []string{"vacationstogo.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 68},
		"cruises_com":        {ID: "cruises_com", Name: "Cruises.com", Domains: []string{"cruises.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 68},
		"outdoorsy":          {ID: "outdoorsy", Name: "Outdoorsy", Domains: []string{"outdoorsy.com"}, SearchTemplate: "https://www.outdoorsy.com/rv-rental/search?keyword={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"rvshare":            {ID: "rvshare", Name: "RVshare", Domains: []string{"rvshare.com"}, SearchTemplate: "https://rvshare.com/search?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"viator":             {ID: "viator", Name: "Viator", Domains: []string{"viator.com"}, SearchTemplate: "https://www.viator.com/searchResults/all?text={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 86},
		"getyourguide":       {ID: "getyourguide", Name: "GetYourGuide", Domains: []string{"getyourguide.com"}, SearchTemplate: "https://www.getyourguide.com/s/?q={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 84},
		"klook":              {ID: "klook", Name: "Klook", Domains: []string{"klook.com"}, SearchTemplate: "https://www.klook.com/search/?query={query}", Modes: []source.Mode{source.ModeHTMLFetch, source.ModeBrowserVerify}, Priority: 78},
		"amadeus":            {ID: "amadeus", Name: "Amadeus", Domains: []string{"amadeus.com"}, Modes: []source.Mode{source.ModeAPI}, Priority: 98},
		"duffel":             {ID: "duffel", Name: "Duffel", Domains: []string{"duffel.com"}, Modes: []source.Mode{source.ModeAPI}, Priority: 96},
		"kiwi":               {ID: "kiwi", Name: "Kiwi", Domains: []string{"kiwi.com"}, SearchTemplate: "https://www.kiwi.com/us/search/results/{query}", Modes: []source.Mode{source.ModeAPI, source.ModeHTMLFetch}, Priority: 88},
		"delta":              {ID: "delta", Name: "Delta", Domains: []string{"delta.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"united":             {ID: "united", Name: "United", Domains: []string{"united.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"american":           {ID: "american", Name: "American Airlines", Domains: []string{"aa.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 70},
		"southwest":          {ID: "southwest", Name: "Southwest", Domains: []string{"southwest.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"jetblue":            {ID: "jetblue", Name: "JetBlue", Domains: []string{"jetblue.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"alaska":             {ID: "alaska", Name: "Alaska Airlines", Domains: []string{"alaskaair.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 65},
		"frontier":           {ID: "frontier", Name: "Frontier", Domains: []string{"flyfrontier.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 55},
		"spirit":             {ID: "spirit", Name: "Spirit", Domains: []string{"spirit.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 55},
		"british_airways":    {ID: "british_airways", Name: "British Airways", Domains: []string{"britishairways.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"lufthansa":          {ID: "lufthansa", Name: "Lufthansa", Domains: []string{"lufthansa.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"air_france":         {ID: "air_france", Name: "Air France", Domains: []string{"airfrance.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
		"klm":                {ID: "klm", Name: "KLM", Domains: []string{"klm.com"}, Modes: []source.Mode{source.ModeBrowserVerify}, Priority: 60},
	}
	collections := map[string]source.Collection{
		"us_domestic_flights":   {ID: "us_domestic_flights", Label: "US domestic flights", Description: "US domestic flight research via API providers first, airline direct as limited verification fallback.", DefaultSources: []string{"amadeus", "duffel", "kiwi", "skyscanner"}, ExtendedSources: []string{"kayak_flights", "google_flights", "momondo_flights", "cheapoair_flights", "delta", "united", "american", "southwest", "jetblue", "alaska"}, MaxSources: 4},
		"international_flights": {ID: "international_flights", Label: "International flights", Description: "International flight research via API/aggregator sources first.", DefaultSources: []string{"amadeus", "duffel", "skyscanner", "kiwi"}, ExtendedSources: []string{"kayak_flights", "google_flights", "momondo_flights", "cheapoair_flights", "united", "delta", "american", "british_airways", "lufthansa", "air_france", "klm"}, MaxSources: 4},
		"budget_flights":        {ID: "budget_flights", Label: "Budget flights", Description: "Budget and low-cost carrier focused flight research with strong fare-trap warnings.", DefaultSources: []string{"kiwi", "skyscanner"}, ExtendedSources: []string{"kayak_flights", "google_flights", "momondo_flights", "cheapoair_flights", "frontier", "spirit", "southwest"}, MaxSources: 4},
		"hotels":                {ID: "hotels", Label: "Hotels", Description: "Hotel and lodging research without booking or payment automation.", DefaultSources: []string{"booking_com", "expedia_hotels", "hotels_com", "priceline_hotels", "travelocity_hotels"}, ExtendedSources: []string{"vrbo"}, MaxSources: 5},
		"camping_outdoors":      {ID: "camping_outdoors", Label: "Camping & outdoors", Description: "Camping, parks, public lands, campgrounds, and outdoor stays research.", DefaultSources: []string{"recreation_gov", "hipcamp", "the_dyrt", "campendium"}, MaxSources: 4},
		"cruises":               {ID: "cruises", Label: "Cruises", Description: "Cruise research with browser-verify-only sources where search is interface-driven.", DefaultSources: []string{"cruise_critic", "vacations_to_go", "cruises_com"}, MaxSources: 3},
		"rv_van_travel":         {ID: "rv_van_travel", Label: "RV & van travel", Description: "RV rentals, camper vans, and road-trip lodging research.", DefaultSources: []string{"outdoorsy", "rvshare"}, MaxSources: 2},
		"tours_activities":      {ID: "tours_activities", Label: "Tours & activities", Description: "Tours, excursions, attractions, and activities research without booking automation.", DefaultSources: []string{"viator", "getyourguide", "klook"}, MaxSources: 3},
	}
	return source.Registry{Sources: sources, Collections: collections}
}
