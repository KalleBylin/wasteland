package main

import (
	"strings"
	"testing"

	"github.com/julianknutsen/wasteland/internal/commons"
)

// Business logic tests for updating moved to internal/sdk/ (sdk_test.go, lifecycle_test.go).

func TestValidateUpdateInputs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		itemType string
		effort   string
		priority int
		wantErr  string
	}{
		{"all empty", "", "", -1, ""},
		{"valid type", "bug", "", -1, ""},
		{"invalid type", "bad", "", -1, "invalid type"},
		{"valid effort", "", "small", -1, ""},
		{"invalid effort", "", "huge", -1, "invalid effort"},
		{"valid priority 0", "", "", 0, ""},
		{"valid priority 4", "", "", 4, ""},
		{"invalid priority too high", "", "", 9, "invalid priority"},
		{"invalid priority negative", "", "", -5, "invalid priority"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateUpdateInputs(tt.itemType, tt.effort, tt.priority)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("validateUpdateInputs() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("validateUpdateInputs() expected error containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestHasUpdateFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		fields *commons.WantedUpdate
		want   bool
	}{
		{"empty", &commons.WantedUpdate{Priority: -1}, false},
		{"title set", &commons.WantedUpdate{Title: "new", Priority: -1}, true},
		{"priority set", &commons.WantedUpdate{Priority: 0}, true},
		{"tags set", &commons.WantedUpdate{Priority: -1, TagsSet: true}, true},
		{"effort set", &commons.WantedUpdate{Priority: -1, EffortLevel: "small"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := hasUpdateFields(tt.fields)
			if got != tt.want {
				t.Errorf("hasUpdateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
