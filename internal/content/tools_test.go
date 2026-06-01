package content

import "testing"

func TestToolMapTargetsKnownCollections(t *testing.T) {
	reg := DefaultRegistry()
	for tool, collection := range ToolMap() {
		if _, ok := reg.Collections[collection]; !ok {
			t.Fatalf("tool %s points to unknown collection %s", tool, collection)
		}
	}
}
