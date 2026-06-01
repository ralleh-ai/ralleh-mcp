package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
)

func TestEvaluateHealthyRegistry(t *testing.T) {
	reg := source.Registry{
		Sources:     map[string]source.Source{"s": {ID: "s"}},
		Collections: map[string]source.Collection{"c": {ID: "c", DefaultSources: []string{"s"}}},
	}
	status := Evaluate("svc", reg, map[string]bool{"canSearch": true})
	if !status.Ready || status.Status != "ok" {
		t.Fatalf("expected healthy status, got %+v", status)
	}
}

func TestReadyzFailsWhenNotReady(t *testing.T) {
	status := Evaluate("svc", source.Registry{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	Handler(status).ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 for unready registry, got %d", w.Code)
	}
}
