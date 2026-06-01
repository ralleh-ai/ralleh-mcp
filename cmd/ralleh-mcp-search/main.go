package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ralleh-ai/ralleh-mcp/internal/content"
	contentsearch "github.com/ralleh-ai/ralleh-mcp/internal/content/search"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/health"
)

func capabilities() map[string]bool {
	return map[string]bool{"canSearchContent": true, "canSummarize": true, "canCrawlArbitraryWebsites": false}
}

func main() {
	healthOnly := flag.Bool("health", false, "print health JSON and exit")
	healthServer := flag.Bool("health-server", false, "serve local-only HTTP health endpoints")
	healthListen := flag.String("health-listen", "127.0.0.1:8624", "health server listen address")
	allowNonLoopback := flag.Bool("allow-non-loopback-health", false, "allow health server to bind outside loopback; requires external firewall/auth controls")
	searchQuery := flag.String("search-query", "", "run deterministic fake content search for smoke/integration testing")
	searchCollection := flag.String("search-collection", "breaking_news", "curated content collection for fake search")
	searchSources := flag.String("search-sources", "", "comma-separated preferred source IDs for fake search")
	rankSources := flag.Bool("rank-sources", false, "print ranked curated source lists and exit")
	rankCollection := flag.String("rank-collection", "", "optional collection ID for --rank-sources")
	flag.Parse()

	reg := content.DefaultRegistry()
	status := health.Evaluate("ralleh-mcp-search", reg, capabilities())
	if *rankSources {
		if *rankCollection != "" {
			ranks, err := reg.RankCollection(*rankCollection)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			writeJSON(map[string]any{"service": "ralleh-mcp-search", "collection": *rankCollection, "rankings": ranks})
			return
		}
		writeJSON(map[string]any{"service": "ralleh-mcp-search", "rankings": reg.RankAll()})
		return
	}
	if *searchQuery != "" {
		resp, err := contentsearch.Search(context.Background(), contentsearch.Request{Query: *searchQuery, Collection: *searchCollection, PreferredSources: splitCSV(*searchSources), BudgetProfile: budget.ProfileFast}, contentsearch.FakeAdapter{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		writeJSON(resp)
		return
	}
	if *healthOnly {
		writeJSON(status)
		if !status.Ready {
			os.Exit(1)
		}
		return
	}
	if *healthServer {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if err := health.ServeLocal(ctx, *healthListen, *allowNonLoopback, status); err != nil && ctx.Err() == nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	writeJSON(map[string]any{"service": "ralleh-mcp-search", "status": "scaffold_ready", "collections": reg.CollectionIDs(), "capabilities": capabilities()})
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
