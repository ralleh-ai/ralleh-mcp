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

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/health"
	"github.com/ralleh-ai/ralleh-mcp/internal/shop"
	shopsearch "github.com/ralleh-ai/ralleh-mcp/internal/shop/search"
)

func capabilities() map[string]bool {
	return map[string]bool{
		"canSearch":        true,
		"canCompare":       false,
		"canVerify":        false,
		"canPurchase":      false,
		"canUseCreditCard": false,
	}
}

func main() {
	healthOnly := flag.Bool("health", false, "print health JSON and exit")
	healthServer := flag.Bool("health-server", false, "serve local-only HTTP health endpoints")
	healthListen := flag.String("health-listen", "127.0.0.1:8621", "health server listen address")
	allowNonLoopback := flag.Bool("allow-non-loopback-health", false, "allow health server to bind outside loopback; requires external firewall/auth controls")
	searchQuery := flag.String("search-query", "", "run deterministic fake shop search for smoke/integration testing")
	searchCollection := flag.String("search-collection", "tools", "curated shop collection for fake search")
	searchSources := flag.String("search-sources", "", "comma-separated preferred source IDs for fake search")
	rankSources := flag.Bool("rank-sources", false, "print ranked curated source lists and exit")
	rankCollection := flag.String("rank-collection", "", "optional collection ID for --rank-sources")
	flag.Parse()

	reg := shop.DefaultRegistry()
	status := health.Evaluate("ralleh-mcp-shop", reg, capabilities())
	if *rankSources {
		if *rankCollection != "" {
			ranks, err := reg.RankCollection(*rankCollection)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			writeJSON(map[string]any{"service": "ralleh-mcp-shop", "collection": *rankCollection, "rankings": ranks})
			return
		}
		writeJSON(map[string]any{"service": "ralleh-mcp-shop", "rankings": reg.RankAll()})
		return
	}
	if *searchQuery != "" {
		resp, err := shopsearch.Search(context.Background(), shopsearch.Request{Query: *searchQuery, Collection: *searchCollection, PreferredSources: splitCSV(*searchSources), BudgetProfile: "fast"}, shopsearch.FakeAdapter{})
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

	b, err := budget.Resolve(budget.ProfileStandard)
	if err != nil {
		panic(err)
	}
	writeJSON(map[string]any{
		"service":       "ralleh-mcp-shop",
		"status":        "scaffold_ready",
		"collections":   reg.CollectionIDs(),
		"defaultBudget": b.Profile,
		"capabilities":  capabilities(),
	})
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
