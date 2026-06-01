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

	"github.com/ralleh-ai/ralleh-mcp/internal/core/health"
	"github.com/ralleh-ai/ralleh-mcp/internal/travel"
	travelsearch "github.com/ralleh-ai/ralleh-mcp/internal/travel/search"
)

func capabilities() map[string]bool {
	return map[string]bool{
		"canSearchFlights":      true,
		"canBook":               false,
		"canUseCreditCard":      false,
		"canEnterPassengerInfo": false,
	}
}

func main() {
	healthOnly := flag.Bool("health", false, "print health JSON and exit")
	healthServer := flag.Bool("health-server", false, "serve local-only HTTP health endpoints")
	healthListen := flag.String("health-listen", "127.0.0.1:8622", "health server listen address")
	allowNonLoopback := flag.Bool("allow-non-loopback-health", false, "allow health server to bind outside loopback; requires external firewall/auth controls")
	flightOrigin := flag.String("flight-origin", "", "run deterministic fake flight search for smoke/integration testing")
	flightDestination := flag.String("flight-destination", "", "destination airport/city for fake flight search")
	flightDepart := flag.String("flight-depart", "", "departure date for fake flight search")
	flightCollection := flag.String("flight-collection", "us_domestic_flights", "curated travel collection for fake flight search")
	flightSources := flag.String("flight-sources", "", "comma-separated preferred source IDs for fake flight search")
	flag.Parse()

	reg := travel.DefaultRegistry()
	status := health.Evaluate("ralleh-mcp-travel", reg, capabilities())
	if *flightOrigin != "" || *flightDestination != "" || *flightDepart != "" {
		resp, err := travelsearch.Search(context.Background(), travelsearch.Request{Origin: *flightOrigin, Destination: *flightDestination, DepartDate: *flightDepart, Collection: *flightCollection, PreferredSources: splitCSV(*flightSources), BudgetProfile: "fast"}, travelsearch.FakeAdapter{})
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

	writeJSON(map[string]any{
		"service":      "ralleh-mcp-travel",
		"status":       "scaffold_ready",
		"collections":  reg.CollectionIDs(),
		"capabilities": capabilities(),
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
