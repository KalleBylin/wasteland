package main

import (
	"testing"
)

// Business logic tests for posting moved to internal/sdk/ (sdk_test.go, lifecycle_test.go).

func TestValidatePostInputs_ValidType(t *testing.T) {
	t.Parallel()
	for _, typ := range []string{"feature", "bug", "design", "rfc", "docs", "inference", ""} {
		if err := validatePostInputs(typ, "medium", 2); err != nil {
			t.Errorf("validatePostInputs(type=%q) unexpected error: %v", typ, err)
		}
	}
}

func TestValidatePostInputs_InvalidType(t *testing.T) {
	t.Parallel()
	err := validatePostInputs("invalid", "medium", 2)
	if err == nil {
		t.Error("validatePostInputs(type=invalid) expected error")
	}
}

func TestValidatePostInputs_InvalidEffort(t *testing.T) {
	t.Parallel()
	err := validatePostInputs("bug", "huge", 2)
	if err == nil {
		t.Error("validatePostInputs(effort=huge) expected error")
	}
}

func TestValidatePostInputs_PriorityBounds(t *testing.T) {
	t.Parallel()
	if err := validatePostInputs("", "medium", -1); err == nil {
		t.Error("validatePostInputs(priority=-1) expected error")
	}
	if err := validatePostInputs("", "medium", 5); err == nil {
		t.Error("validatePostInputs(priority=5) expected error")
	}
	for _, p := range []int{0, 1, 2, 3, 4} {
		if err := validatePostInputs("", "medium", p); err != nil {
			t.Errorf("validatePostInputs(priority=%d) unexpected error: %v", p, err)
		}
	}
}
