// Package baseline provides functionality to save and load a reference
// snapshot of diff results, enabling comparison against a known-good state.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot represents a saved baseline of diff results.
type Snapshot struct {
	CreatedAt time.Time   `json:"created_at"`
	FileA     string      `json:"file_a"`
	FileB     string      `json:"file_b"`
	Results   []diff.Result `json:"results"`
}

// Save writes the given results to a baseline file at the specified path.
func Save(path, fileA, fileB string, results []diff.Result) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		FileA:     fileA,
		FileB:     fileB,
		Results:   results,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write file %q: %w", path, err)
	}

	return nil
}

// Load reads a baseline snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read file %q: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal snapshot: %w", err)
	}

	return &snap, nil
}

// NewIssues returns results that are present in current but were not in the
// baseline snapshot, indicating regressions introduced since the snapshot.
func NewIssues(snap *Snapshot, current []diff.Result) []diff.Result {
	baselineKeys := make(map[string]struct{}, len(snap.Results))
	for _, r := range snap.Results {
		baselineKeys[r.Key] = struct{}{}
	}

	var novel []diff.Result
	for _, r := range current {
		if _, exists := baselineKeys[r.Key]; !exists {
			novel = append(novel, r)
		}
	}

	return novel
}
