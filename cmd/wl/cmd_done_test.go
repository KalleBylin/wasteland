package main

import (
	"strings"
	"testing"

	"github.com/gastownhall/wasteland/internal/commons"
)

// Business logic tests for submitting completions moved to internal/sdk/.

func TestGenerateCompletionID_Format(t *testing.T) {
	t.Parallel()
	id := commons.GeneratePrefixedID("c", "w-abc123", "my-rig")
	if !strings.HasPrefix(id, "c-") {
		t.Errorf("GeneratePrefixedID(c) = %q, want prefix 'c-'", id)
	}
	// "c-" + 16 hex chars = 18 chars total
	if len(id) != 18 {
		t.Errorf("GeneratePrefixedID(c) length = %d, want 18", len(id))
	}
	hexPart := id[2:]
	for _, c := range hexPart {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("GeneratePrefixedID(c) contains non-hex char %q in %q", string(c), id)
		}
	}
}

func TestGenerateCompletionID_DeterministicInputs(t *testing.T) {
	t.Parallel()
	id1 := commons.GeneratePrefixedID("c", "w-abc", "rig-1")
	id2 := commons.GeneratePrefixedID("c", "w-def", "rig-1")
	id3 := commons.GeneratePrefixedID("c", "w-abc", "rig-2")

	if id1 == id2 {
		t.Errorf("same ID for different wantedIDs: %s", id1)
	}
	if id1 == id3 {
		t.Errorf("same ID for different rigHandles: %s", id1)
	}
}
