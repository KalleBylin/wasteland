package main

import (
	"strings"
	"testing"

	"github.com/gastownhall/wasteland/internal/commons"
)

// Business logic tests for accepting moved to internal/sdk/ (sdk_test.go, lifecycle_test.go).

func TestGenerateStampID_Format(t *testing.T) {
	t.Parallel()
	id := commons.GeneratePrefixedID("s", "w-abc123", "my-rig")
	if !strings.HasPrefix(id, "s-") {
		t.Errorf("GeneratePrefixedID(s) = %q, want prefix 's-'", id)
	}
	// "s-" + 16 hex chars = 18 chars total
	if len(id) != 18 {
		t.Errorf("GeneratePrefixedID(s) length = %d, want 18", len(id))
	}
	hexPart := id[2:]
	for _, c := range hexPart {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("GeneratePrefixedID(s) contains non-hex char %q in %q", string(c), id)
		}
	}
}

func TestValidateAcceptInputs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		quality     int
		reliability int
		severity    string
		wantErr     string
	}{
		{"valid", 3, 4, "leaf", ""},
		{"quality too low", 0, 3, "leaf", "invalid quality"},
		{"quality too high", 6, 3, "leaf", "invalid quality"},
		{"reliability too low", 3, 0, "leaf", "invalid reliability"},
		{"reliability too high", 3, 6, "leaf", "invalid reliability"},
		{"bad severity", 3, 3, "bad", "invalid severity"},
		{"valid branch", 5, 5, "branch", ""},
		{"valid root", 1, 1, "root", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateAcceptInputs(tt.quality, tt.reliability, tt.severity)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("validateAcceptInputs() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("validateAcceptInputs() expected error containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}
