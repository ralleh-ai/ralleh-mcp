package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/mcpstdio"
)

func TestSearchMCPToolsOverStdio(t *testing.T) {
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search_content","arguments":{"query":"ai chips","collection":"technology","preferredSources":["hacker_news","random_blog"],"budgetProfile":"fast"}}}`,
	}, "\n") + "\n"

	var out strings.Builder
	if err := mcpstdio.Serve(context.Background(), strings.NewReader(input), &out, "ralleh-mcp-search", "test", Tools()); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 JSON-RPC responses, got %d: %s", len(lines), out.String())
	}
	assertJSONRPCID(t, lines[0], float64(1))
	list := assertJSONRPCID(t, lines[1], float64(2))
	tools := list["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(tools))
	}
	call := assertJSONRPCID(t, lines[2], float64(3))
	result := call["result"].(map[string]any)
	content := result["content"].([]any)
	text := content[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, `"status": "ok"`) || !strings.Contains(text, "random_blog") {
		t.Fatalf("expected successful search payload with rejected source, got %s", text)
	}
}

func assertJSONRPCID(t *testing.T, line string, want float64) map[string]any {
	t.Helper()
	var msg map[string]any
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		t.Fatalf("invalid JSON response %q: %v", line, err)
	}
	if msg["jsonrpc"] != "2.0" || msg["id"] != want {
		t.Fatalf("unexpected response id/protocol: %+v", msg)
	}
	if _, hasErr := msg["error"]; hasErr {
		t.Fatalf("unexpected JSON-RPC error: %+v", msg)
	}
	return msg
}
