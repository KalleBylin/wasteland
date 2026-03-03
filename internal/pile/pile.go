// Package pile provides a read-only DoltHub client for hop/the-pile.
//
// This is intentionally separate from the SDK/backend stack, which is built
// around joined wastelands with fork semantics. The pile is a global read-only
// database of seeded developer profiles.
package pile

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/julianknutsen/wasteland/internal/backend"
)

// PileClient is a read-only DoltHub API client for hop/the-pile.
type PileClient struct {
	org    string
	db     string
	branch string
	token  string
	client *http.Client
}

// New creates a PileClient targeting the given org/db on DoltHub.
// Token is optional for public databases.
func New(token, org, db string) *PileClient {
	return &PileClient{
		org:    org,
		db:     db,
		branch: "main",
		token:  token,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewDefault creates a PileClient for hop/the-pile with no auth token.
func NewDefault() *PileClient {
	return New("", "hop", "the-pile")
}

// Query runs a read-only SQL query and returns the result as CSV.
func (p *PileClient) Query(sql string) (string, error) {
	body, err := p.queryRaw(sql)
	if err != nil {
		return "", err
	}
	return backend.JSONToCSV(body)
}

// queryRaw runs a SQL query and returns the raw DoltHub JSON response body.
func (p *PileClient) queryRaw(sql string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/%s/%s/%s?q=%s",
		backend.DoltHubAPIBase, p.org, p.db,
		url.PathEscape(p.branch), url.QueryEscape(sql))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if p.token != "" {
		req.Header.Set("Authorization", "token "+p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pile query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoltHub API returned %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// QueryRows runs a SQL query and returns parsed JSON rows.
func (p *PileClient) QueryRows(sql string) ([]map[string]any, error) {
	body, err := p.queryRaw(sql)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status  string            `json:"query_execution_status"`
		Message string            `json:"query_execution_message"`
		Rows    []json.RawMessage `json:"rows"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if resp.Status == "Error" {
		return nil, fmt.Errorf("query error: %s", resp.Message)
	}

	rows := make([]map[string]any, 0, len(resp.Rows))
	for _, raw := range resp.Rows {
		var row map[string]any
		if err := json.Unmarshal(raw, &row); err != nil {
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}
