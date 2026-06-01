package budget

import (
	"fmt"
	"time"
)

// Profile is a named request budget. Callers may ask for a profile, but the
// service owns the actual limits and clamps everything to safe values.
type Profile string

const (
	ProfileFast     Profile = "fast"
	ProfileStandard Profile = "standard"
	ProfileDeep     Profile = "deep"
)

// Budget captures hard execution limits for one MCP call.
type Budget struct {
	Profile             Profile
	GlobalTimeout       time.Duration
	PerSourceTimeout    time.Duration
	MaxSources          int
	MaxConcurrency      int
	MaxResultsPerSource int
	MaxResponseBytes    int64
	MaxRetriesPerSource int
}

// Resolve returns the effective server-owned budget for a profile.
func Resolve(profile Profile) (Budget, error) {
	switch profile {
	case "", ProfileStandard:
		return Budget{
			Profile:             ProfileStandard,
			GlobalTimeout:       12 * time.Second,
			PerSourceTimeout:    3500 * time.Millisecond,
			MaxSources:          5,
			MaxConcurrency:      4,
			MaxResultsPerSource: 5,
			MaxResponseBytes:    750_000,
			MaxRetriesPerSource: 1,
		}, nil
	case ProfileFast:
		return Budget{
			Profile:             ProfileFast,
			GlobalTimeout:       6 * time.Second,
			PerSourceTimeout:    2500 * time.Millisecond,
			MaxSources:          3,
			MaxConcurrency:      3,
			MaxResultsPerSource: 3,
			MaxResponseBytes:    500_000,
			MaxRetriesPerSource: 0,
		}, nil
	case ProfileDeep:
		return Budget{
			Profile:             ProfileDeep,
			GlobalTimeout:       20 * time.Second,
			PerSourceTimeout:    5 * time.Second,
			MaxSources:          8,
			MaxConcurrency:      4,
			MaxResultsPerSource: 8,
			MaxResponseBytes:    1_000_000,
			MaxRetriesPerSource: 1,
		}, nil
	default:
		return Budget{}, fmt.Errorf("unknown budget profile %q", profile)
	}
}
