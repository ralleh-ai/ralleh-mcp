package budget

import "testing"

func TestResolveClampedProfiles(t *testing.T) {
	for _, profile := range []Profile{ProfileFast, ProfileStandard, ProfileDeep, ""} {
		b, err := Resolve(profile)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", profile, err)
		}
		if b.GlobalTimeout <= 0 || b.PerSourceTimeout <= 0 {
			t.Fatalf("profile %q returned invalid timeouts: %+v", profile, b)
		}
		if b.MaxConcurrency <= 0 || b.MaxConcurrency > 4 {
			t.Fatalf("profile %q returned unsafe concurrency: %+v", profile, b)
		}
		if b.MaxResponseBytes <= 0 || b.MaxResponseBytes > 1_000_000 {
			t.Fatalf("profile %q returned unsafe response cap: %+v", profile, b)
		}
	}
}

func TestResolveRejectsUnknownProfile(t *testing.T) {
	if _, err := Resolve("crawl_the_world"); err == nil {
		t.Fatal("expected unknown profile to be rejected")
	}
}
