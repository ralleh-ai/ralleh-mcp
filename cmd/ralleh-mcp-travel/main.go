package main

import (
	"encoding/json"
	"os"

	"github.com/ralleh-ai/ralleh-mcp/internal/travel"
)

func main() {
	reg := travel.DefaultRegistry()
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"service":     "ralleh-mcp-travel",
		"status":      "scaffold_ready",
		"collections": reg.CollectionIDs(),
		"capabilities": map[string]bool{
			"canSearchFlights":      true,
			"canBook":               false,
			"canUseCreditCard":      false,
			"canEnterPassengerInfo": false,
		},
	})
}
