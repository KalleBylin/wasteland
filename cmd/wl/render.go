package main

import (
	"fmt"
	"io"

	"github.com/gastownhall/wasteland/internal/sdk"
	"github.com/gastownhall/wasteland/internal/style"
)

// renderMutationResult writes a consistent summary for SDK mutation results.
//
//	verb:    past-tense action word, e.g. "Claimed", "Unclaimed", "Deleted"
//	wantedID: the item ID
//	result:  the SDK mutation result
//	extras:  additional "key: value" lines to render between the header and hint
func renderMutationResult(w io.Writer, verb, wantedID string, result *sdk.MutationResult, extras ...string) {
	fmt.Fprintf(w, "%s %s %s\n", style.Bold.Render("✓"), verb, wantedID)

	if result.Detail != nil && result.Detail.Item != nil {
		fmt.Fprintf(w, "  Title: %s\n", result.Detail.Item.Title)
		if result.Detail.Item.Status != "" {
			fmt.Fprintf(w, "  Status: %s\n", result.Detail.Item.Status)
		}
	}

	for _, extra := range extras {
		fmt.Fprintf(w, "  %s\n", extra)
	}

	if result.Branch != "" {
		fmt.Fprintf(w, "  Branch: %s\n", result.Branch)
	}
	if result.Detail != nil && result.Detail.BranchURL != "" {
		fmt.Fprintf(w, "  Branch URL: %s\n", result.Detail.BranchURL)
	}
	if result.Detail != nil && result.Detail.PRURL != "" {
		fmt.Fprintf(w, "  PR: %s\n", result.Detail.PRURL)
	}
	if result.Detail != nil && result.Detail.Delta != "" {
		fmt.Fprintf(w, "  Delta: %s\n", result.Detail.Delta)
	}

	if result.Hint != "" {
		fmt.Fprintf(w, "\n  %s\n", style.Dim.Render(result.Hint))
	}
}
