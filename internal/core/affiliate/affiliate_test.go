package affiliate

import "testing"

func TestApplyQueryParam(t *testing.T) {
	rule := Rule{SourceID: "amazon_api", AllowedDomains: []string{"amazon.com"}, Param: "tag", Value: "ralleh-20", Enabled: true}
	res, err := ApplyQueryParam("https://www.amazon.com/dp/B00TEST?psc=1", rule)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied || res.PresentedURL == res.CanonicalURL {
		t.Fatalf("expected affiliate tag, got %+v", res)
	}
}

func TestApplyQueryParamRejectsUnsafeURLs(t *testing.T) {
	rule := Rule{SourceID: "amazon_api", AllowedDomains: []string{"amazon.com"}, Param: "tag", Value: "ralleh-20", Enabled: true}
	for _, raw := range []string{"http://www.amazon.com/dp/B00TEST", "https://evil.example/dp/B00TEST", "https://www.amazon.com/checkout"} {
		res, err := ApplyQueryParam(raw, rule)
		if err != nil {
			t.Fatal(err)
		}
		if res.Applied {
			t.Fatalf("expected unsafe URL not to be tagged: %+v", res)
		}
	}
}
