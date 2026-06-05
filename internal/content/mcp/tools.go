package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/ralleh-ai/ralleh-mcp/internal/content"
	contentsearch "github.com/ralleh-ai/ralleh-mcp/internal/content/search"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/mcpstdio"
)

func Tools() []mcpstdio.Tool {
	return []mcpstdio.Tool{
		{
			Name:        "list_collections",
			Description: "List curated content/news/research collections available to Ralleh search.",
			InputSchema: objectSchema(map[string]any{}, nil),
			Handler: func(_ context.Context, _ map[string]any) (any, error) {
				reg := content.DefaultRegistry()
				items := make([]map[string]any, 0, len(reg.Collections))
				for _, id := range reg.CollectionIDs() {
					c := reg.Collections[id]
					items = append(items, map[string]any{"id": c.ID, "label": c.Label, "description": c.Description, "maxSources": c.MaxSources})
				}
				return map[string]any{"collections": items}, nil
			},
		},
		{
			Name:        "rank_sources",
			Description: "Rank curated sources for one content collection, or all collections when collection is omitted.",
			InputSchema: objectSchema(map[string]any{"collection": map[string]any{"type": "string", "description": "Optional collection ID."}}, nil),
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				reg := content.DefaultRegistry()
				collection := strings.TrimSpace(stringArg(args, "collection"))
				if collection == "" {
					return map[string]any{"rankings": reg.RankAll()}, nil
				}
				rankings, err := reg.RankCollection(collection)
				if err != nil {
					return nil, err
				}
				return map[string]any{"collection": collection, "rankings": rankings}, nil
			},
		},
		{
			Name:        "search_content",
			Description: "Run bounded curated content search using known collection/source IDs only. Does not crawl arbitrary websites.",
			InputSchema: objectSchema(map[string]any{
				"query":            map[string]any{"type": "string", "description": "Search query."},
				"collection":       map[string]any{"type": "string", "description": "Curated collection ID, for example technology, breaking_news, research, stocks_markets."},
				"preferredSources": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Optional preferred known source IDs."},
				"budgetProfile":    map[string]any{"type": "string", "enum": []string{"fast", "standard", "deep"}, "description": "Optional budget profile; defaults to fast."},
			}, []string{"query", "collection"}),
			Handler: func(ctx context.Context, args map[string]any) (any, error) {
				query := strings.TrimSpace(stringArg(args, "query"))
				collection := strings.TrimSpace(stringArg(args, "collection"))
				if query == "" {
					return nil, fmt.Errorf("query is required")
				}
				if collection == "" {
					return nil, fmt.Errorf("collection is required")
				}
				profile := budget.Profile(strings.TrimSpace(stringArg(args, "budgetProfile")))
				if profile == "" {
					profile = budget.ProfileFast
				}
				return contentsearch.Search(ctx, contentsearch.Request{Query: query, Collection: collection, PreferredSources: stringSliceArg(args, "preferredSources"), BudgetProfile: profile}, contentsearch.FakeAdapter{})
			},
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	out := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		out["required"] = required
	}
	return out
}

func stringArg(args map[string]any, key string) string {
	if args == nil {
		return ""
	}
	v, _ := args[key].(string)
	return v
}

func stringSliceArg(args map[string]any, key string) []string {
	if args == nil {
		return nil
	}
	raw, ok := args[key]
	if !ok || raw == nil {
		return nil
	}
	values, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		s, ok := v.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
