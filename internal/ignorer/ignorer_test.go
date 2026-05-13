package ignorer_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.envignore")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Kind: diff.Missing},
		{Key: "DB_PASSWORD", Kind: diff.Mismatch},
		{Key: "SECRET_KEY", Kind: diff.Missing},
		{Key: "APP_ENV", Kind: diff.Mismatch},
		{Key: "INTERNAL_TOKEN", Kind: diff.Missing},
	}
}

func TestNew_Matches_Exact(t *testing.T) {
	il := ignorer.New([]string{"DB_HOST", "SECRET_KEY"})
	if !il.Matches("DB_HOST") {
		t.Error("expected DB_HOST to match")
	}
	if il.Matches("APP_ENV") {
		t.Error("expected APP_ENV not to match")
	}
}

func TestNew_Matches_Prefix(t *testing.T) {
	il := ignorer.New([]string{"DB_*", "INTERNAL_*"})
	if !il.Matches("DB_HOST") {
		t.Error("expected DB_HOST to match prefix DB_*")
	}
	if !il.Matches("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to match prefix DB_*")
	}
	if !il.Matches("INTERNAL_TOKEN") {
		t.Error("expected INTERNAL_TOKEN to match prefix INTERNAL_*")
	}
	if il.Matches("APP_ENV") {
		t.Error("expected APP_ENV not to match")
	}
}

func TestApply_FiltersResults(t *testing.T) {
	il := ignorer.New([]string{"DB_*", "SECRET_KEY"})
	results := sampleResults()
	got := il.Apply(results)
	if len(got) != 2 {
		t.Fatalf("expected 2 results after filtering, got %d", len(got))
	}
	for _, r := range got {
		if r.Key == "DB_HOST" || r.Key == "DB_PASSWORD" || r.Key == "SECRET_KEY" {
			t.Errorf("key %q should have been filtered out", r.Key)
		}
	}
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempIgnore(t, "# ignore secrets\nSECRET_KEY\nDB_*\n\nINTERNAL_TOKEN\n")
	il, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if !il.Matches("SECRET_KEY") {
		t.Error("expected SECRET_KEY to match")
	}
	if !il.Matches("DB_HOST") {
		t.Error("expected DB_HOST to match via prefix")
	}
	if !il.Matches("INTERNAL_TOKEN") {
		t.Error("expected INTERNAL_TOKEN to match")
	}
	if il.Matches("APP_ENV") {
		t.Error("expected APP_ENV not to match")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := ignorer.LoadFile("/nonexistent/path/.envignore")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
