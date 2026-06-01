package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/runtime"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
)

// Status is the machine-readable health payload used by CLI health checks and
// local-only HTTP health endpoints.
type Status struct {
	Service      string          `json:"service"`
	Status       string          `json:"status"`
	Ready        bool            `json:"ready"`
	Version      string          `json:"version"`
	Commit       string          `json:"commit"`
	BuildDate    string          `json:"buildDate"`
	CheckedAt    time.Time       `json:"checkedAt"`
	Collections  []string        `json:"collections"`
	Budgets      []string        `json:"budgets"`
	Issues       []string        `json:"issues,omitempty"`
	Capabilities map[string]bool `json:"capabilities"`
}

// Evaluate verifies static service readiness without touching external sites.
func Evaluate(service string, registry source.Registry, capabilities map[string]bool) Status {
	issues := []string{}
	collections := registry.CollectionIDs()
	if len(collections) == 0 {
		issues = append(issues, "no_collections_configured")
	}
	if len(registry.Sources) == 0 {
		issues = append(issues, "no_sources_configured")
	}
	for _, profile := range []budget.Profile{budget.ProfileFast, budget.ProfileStandard, budget.ProfileDeep} {
		if _, err := budget.Resolve(profile); err != nil {
			issues = append(issues, "budget_profile_invalid:"+string(profile))
		}
	}
	status := "ok"
	ready := true
	if len(issues) > 0 {
		status = "degraded"
		ready = false
	}
	return Status{
		Service:      service,
		Status:       status,
		Ready:        ready,
		Version:      runtime.Version,
		Commit:       runtime.Commit,
		BuildDate:    runtime.Date,
		CheckedAt:    time.Now().UTC(),
		Collections:  collections,
		Budgets:      []string{string(budget.ProfileFast), string(budget.ProfileStandard), string(budget.ProfileDeep)},
		Issues:       issues,
		Capabilities: capabilities,
	}
}

// Handler returns a local health handler for /healthz, /readyz, and /version.
func Handler(status Status) http.Handler {
	mux := http.NewServeMux()
	write := func(w http.ResponseWriter, code int, payload any) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(payload)
	}
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		write(w, http.StatusOK, status)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !status.Ready {
			write(w, http.StatusServiceUnavailable, status)
			return
		}
		write(w, http.StatusOK, status)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		write(w, http.StatusOK, map[string]string{"version": status.Version, "commit": status.Commit, "buildDate": status.BuildDate})
	})
	return mux
}
