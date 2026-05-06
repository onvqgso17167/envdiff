package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "APP_HOST", Status: diff.MissingInB},
	{Key: "APP_PORT", Status: diff.Mismatch, ValueA: "8080", ValueB: "9090"},
	{Key: "DB_HOST", Status: diff.MissingInA},
	{Key: "DB_PASS", Status: diff.Mismatch, ValueA: "secret", ValueB: "other"},
}

func TestApply_NoFilter(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{})
	if len(out) != len(sampleResults) {
		t.Fatalf("expected %d results, got %d", len(sampleResults), len(out))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{OnlyMissing: true})
	for _, r := range out {
		if r.Status != diff.MissingInA && r.Status != diff.MissingInB {
			t.Errorf("unexpected status %q for key %q", r.Status, r.Key)
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 missing results, got %d", len(out))
	}
}

func TestApply_OnlyMismatch(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{OnlyMismatch: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 mismatch results, got %d", len(out))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Prefix: "DB_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 DB_ results, got %d", len(out))
	}
	for _, r := range out {
		if r.Key[:3] != "DB_" {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestApply_KeysFilter(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Keys: []string{"APP_PORT", "DB_HOST"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_CombinedPrefixAndMismatch(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Prefix: "DB_", OnlyMismatch: true})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %q", out[0].Key)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := filter.Apply(nil, filter.Options{OnlyMissing: true})
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d", len(out))
	}
}
