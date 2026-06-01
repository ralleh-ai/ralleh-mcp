package travel

import "testing"

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
