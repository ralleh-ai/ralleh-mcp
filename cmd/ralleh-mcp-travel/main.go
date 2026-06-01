package main

import (
	"encoding/json"
	"os"
)

func main() {
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"service": "ralleh-mcp-travel",
		"status":  "scaffold_ready",
		"capabilities": map[string]bool{
			"canSearchFlights":      true,
			"canBook":               false,
			"canUseCreditCard":      false,
			"canEnterPassengerInfo": false,
		},
	})
}
