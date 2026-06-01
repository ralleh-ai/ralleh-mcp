package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/budget"
	"github.com/ralleh-ai/ralleh-mcp/internal/shop"
)

func main() {
	reg := shop.DefaultRegistry()
	b, err := budget.Resolve(budget.ProfileStandard)
	if err != nil {
		panic(err)
	}
	out := map[string]any{
		"service":       "ralleh-mcp-shop",
		"status":        "scaffold_ready",
		"collections":   reg.CollectionIDs(),
		"defaultBudget": b.Profile,
		"capabilities": map[string]bool{
			"canSearch":        true,
			"canCompare":       false,
			"canVerify":        false,
			"canPurchase":      false,
			"canUseCreditCard": false,
		},
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
