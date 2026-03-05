package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gastownhall/wasteland/internal/commons"
	"github.com/gastownhall/wasteland/internal/inference"
	"github.com/gastownhall/wasteland/internal/sdk"
)

// Tests that modify inference.OllamaURL must not use t.Parallel().

// --- executeInferVerify tests (new signature: takes *sdk.DetailResult) ---

func TestExecuteInferVerify_Match(t *testing.T) {
	output := "The answer is 2."
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			Response string `json:"response"`
		}{Response: output})
	}))
	defer srv.Close()

	old := inference.OllamaURL
	inference.OllamaURL = srv.URL
	defer func() { inference.OllamaURL = old }()

	job := &inference.Job{Prompt: "what is 1+1", Model: "llama3.2:1b", Seed: 42}
	desc, _ := inference.EncodeJob(job)
	result := &inference.Result{
		Output:     output,
		OutputHash: inference.Hash(output),
		Model:      "llama3.2:1b",
		Seed:       42,
	}
	evidence, _ := inference.EncodeResult(result)

	detail := &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-verify1",
			Title:       "infer: test",
			Description: desc,
			Type:        "inference",
		},
		Completion: &commons.CompletionRecord{
			ID:          "c-verify1",
			WantedID:    "w-verify1",
			CompletedBy: "bob",
			Evidence:    evidence,
		},
	}

	vr, err := executeInferVerify(detail, "w-verify1")
	if err != nil {
		t.Fatalf("executeInferVerify() error: %v", err)
	}
	if !vr.Match {
		t.Errorf("Match = false, want true")
	}
}

func TestExecuteInferVerify_Mismatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			Response string `json:"response"`
		}{Response: "different output"})
	}))
	defer srv.Close()

	old := inference.OllamaURL
	inference.OllamaURL = srv.URL
	defer func() { inference.OllamaURL = old }()

	job := &inference.Job{Prompt: "test", Model: "m", Seed: 1}
	desc, _ := inference.EncodeJob(job)
	result := &inference.Result{
		Output:     "original output",
		OutputHash: inference.Hash("original output"),
		Model:      "m",
		Seed:       1,
	}
	evidence, _ := inference.EncodeResult(result)

	detail := &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-verify2",
			Title:       "infer: test",
			Description: desc,
			Type:        "inference",
		},
		Completion: &commons.CompletionRecord{
			ID:          "c-verify2",
			WantedID:    "w-verify2",
			CompletedBy: "bob",
			Evidence:    evidence,
		},
	}

	vr, err := executeInferVerify(detail, "w-verify2")
	if err != nil {
		t.Fatalf("executeInferVerify() error: %v", err)
	}
	if vr.Match {
		t.Error("Match = true, want false")
	}
}

func TestExecuteInferVerify_WrongType(t *testing.T) {
	t.Parallel()
	detail := &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:    "w-wrongtype",
			Title: "Not inference",
			Type:  "bug",
		},
	}

	_, err := executeInferVerify(detail, "w-wrongtype")
	if err == nil {
		t.Fatal("expected error for wrong type")
	}
	if !strings.Contains(err.Error(), "inference") {
		t.Errorf("error = %q, want to mention 'inference'", err.Error())
	}
}

func TestExecuteInferVerify_NoCompletion(t *testing.T) {
	t.Parallel()
	job := &inference.Job{Prompt: "test", Model: "m", Seed: 1}
	desc, _ := inference.EncodeJob(job)

	detail := &sdk.DetailResult{
		Item: &commons.WantedItem{
			ID:          "w-nocomp",
			Title:       "No completion",
			Description: desc,
			Type:        "inference",
		},
	}

	_, err := executeInferVerify(detail, "w-nocomp")
	if err == nil {
		t.Fatal("expected error for missing completion")
	}
}
