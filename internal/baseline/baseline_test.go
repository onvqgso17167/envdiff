package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/baseline"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Kind: diff.Missing, ValueA: "production", ValueB: ""},
		{Key: "DB_HOST", Kind: diff.Mismatch, ValueA: "localhost", ValueB: "db.prod"},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	results := sampleResults()
	if err := baseline.Save(path, "dev.env", "prod.env", results); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if snap.FileA != "dev.env" || snap.FileB != "prod.env" {
		t.Errorf("file names mismatch: got %q %q", snap.FileA, snap.FileB)
	}
	if len(snap.Results) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.Results[0].Key != "APP_ENV" {
		t.Errorf("unexpected key: %s", snap.Results[0].Key)
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestNewIssues_DetectsRegression(t *testing.T) {
	snap := &baseline.Snapshot{
		CreatedAt: time.Now(),
		Results:   sampleResults(),
	}

	current := append(sampleResults(), diff.Result{
		Key:  "NEW_KEY",
		Kind: diff.Missing,
	})

	novel := baseline.NewIssues(snap, current)
	if len(novel) != 1 {
		t.Fatalf("expected 1 new issue, got %d", len(novel))
	}
	if novel[0].Key != "NEW_KEY" {
		t.Errorf("unexpected key: %s", novel[0].Key)
	}
}

func TestNewIssues_NoRegression(t *testing.T) {
	snap := &baseline.Snapshot{
		CreatedAt: time.Now(),
		Results:   sampleResults(),
	}

	novel := baseline.NewIssues(snap, sampleResults())
	if len(novel) != 0 {
		t.Errorf("expected no new issues, got %d", len(novel))
	}
}
