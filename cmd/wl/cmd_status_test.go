package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/julianknutsen/wasteland/internal/commons"
	"github.com/julianknutsen/wasteland/internal/sdk"
)

// Business logic tests for status querying moved to internal/sdk/ (reads.go).

func TestRenderDetailStatus_Open(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderDetailStatus(&buf, &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-abc123",
			Title:       "Fix the login bug",
			Status:      "open",
			Type:        "bug",
			Priority:    1,
			Project:     "gastown",
			EffortLevel: "medium",
			PostedBy:    "poster-rig",
			Tags:        []string{"go", "auth"},
			CreatedAt:   "2026-02-20 14:30:05",
			UpdatedAt:   "2026-02-20 14:30:05",
			Description: "The login page crashes.",
		},
	})

	out := buf.String()
	if !strings.Contains(out, "w-abc123") {
		t.Errorf("output missing wanted ID")
	}
	if !strings.Contains(out, "Fix the login bug") {
		t.Errorf("output missing title")
	}
	if !strings.Contains(out, "open") {
		t.Errorf("output missing status")
	}
	if !strings.Contains(out, "bug") {
		t.Errorf("output missing type")
	}
	if !strings.Contains(out, "P1") {
		t.Errorf("output missing priority")
	}
	if !strings.Contains(out, "gastown") {
		t.Errorf("output missing project")
	}
	if !strings.Contains(out, "poster-rig") {
		t.Errorf("output missing posted by")
	}
	if !strings.Contains(out, "go, auth") {
		t.Errorf("output missing tags")
	}
	if !strings.Contains(out, "The login page crashes.") {
		t.Errorf("output missing description")
	}
	// Should NOT contain completion or stamp sections
	if strings.Contains(out, "Completion:") {
		t.Errorf("output should not contain Completion section for open item")
	}
	if strings.Contains(out, "Stamp:") {
		t.Errorf("output should not contain Stamp section for open item")
	}
}

func TestRenderDetailStatus_Completed(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderDetailStatus(&buf, &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-abc123",
			Title:       "Fix the login bug",
			Status:      "completed",
			Type:        "bug",
			Priority:    1,
			Project:     "gastown",
			EffortLevel: "medium",
			PostedBy:    "poster-rig",
			Tags:        []string{"go", "auth"},
			ClaimedBy:   "worker-rig",
			CreatedAt:   "2026-02-20 14:30:05",
			UpdatedAt:   "2026-02-23 09:15:00",
		},
		Completion: &commons.CompletionRecord{
			ID:          "c-abc123def456ab",
			WantedID:    "w-abc123",
			CompletedBy: "worker-rig",
			Evidence:    "https://github.com/org/repo/pull/123",
		},
		Stamp: &commons.Stamp{
			ID:          "s-abc123def456ab",
			Author:      "reviewer-rig",
			Subject:     "worker-rig",
			Quality:     4,
			Reliability: 3,
			Severity:    "leaf",
			SkillTags:   []string{"go", "auth"},
			Message:     "solid work",
		},
	})

	out := buf.String()
	if !strings.Contains(out, "completed") {
		t.Errorf("output missing status")
	}
	if !strings.Contains(out, "Claimed by:") {
		t.Errorf("output missing claimed by")
	}
	if !strings.Contains(out, "c-abc123def456ab") {
		t.Errorf("output missing completion ID")
	}
	if !strings.Contains(out, "https://github.com/org/repo/pull/123") {
		t.Errorf("output missing evidence")
	}
	if !strings.Contains(out, "s-abc123def456ab") {
		t.Errorf("output missing stamp ID")
	}
	if !strings.Contains(out, "Quality: 4") {
		t.Errorf("output missing quality")
	}
	if !strings.Contains(out, "Reliability: 3") {
		t.Errorf("output missing reliability")
	}
	if !strings.Contains(out, "Severity: leaf") {
		t.Errorf("output missing severity")
	}
	if !strings.Contains(out, "go, auth") {
		t.Errorf("output missing skill tags")
	}
	if !strings.Contains(out, "Accepted by: reviewer-rig") {
		t.Errorf("output missing accepted by")
	}
	if !strings.Contains(out, "Message:     solid work") {
		t.Errorf("output missing message")
	}
}

func TestRenderDetailStatus_PRMode(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderDetailStatus(&buf, &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-pr123",
			Title:       "PR mode test",
			Status:      "claimed",
			Type:        "feature",
			Priority:    2,
			EffortLevel: "medium",
			PostedBy:    "poster-rig",
			ClaimedBy:   "worker-rig",
		},
		Branch:     "wl/worker-rig/w-pr123",
		BranchURL:  "https://dolthub.com/repo/data/wl%2Fworker-rig%2Fw-pr123",
		MainStatus: "open",
		Delta:      "open → claimed",
		PRURL:      "https://dolthub.com/repo/pulls/42",
	})

	out := buf.String()
	if !strings.Contains(out, "wl/worker-rig/w-pr123") {
		t.Errorf("output missing branch")
	}
	if !strings.Contains(out, "open → claimed") {
		t.Errorf("output missing delta")
	}
	if !strings.Contains(out, "https://dolthub.com/repo/pulls/42") {
		t.Errorf("output missing PR URL")
	}
}
