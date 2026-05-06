package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Kind: "mismatch", ValueA: "x", ValueB: "y"},
		{Key: "ALPHA", Kind: "missing_in_b", ValueA: "1", ValueB: ""},
		{Key: "MANGO", Kind: "missing_in_a", ValueA: "", ValueB: "2"},
		{Key: "BETA", Kind: "mismatch", ValueA: "a", ValueB: "b"},
	}
}

func TestSort_ByKey(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.SortByKey)

	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByKind(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.SortByKind)

	// missing_in_a < missing_in_b < mismatch alphabetically
	if sorted[0].Kind != "missing_in_a" {
		t.Errorf("expected first kind to be missing_in_a, got %q", sorted[0].Kind)
	}
	if sorted[1].Kind != "missing_in_b" {
		t.Errorf("expected second kind to be missing_in_b, got %q", sorted[1].Kind)
	}
}

func TestSort_BySeverity(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.SortBySeverity)

	// severity: missing_in_a=1, missing_in_b=2, mismatch=3
	if sorted[0].Kind != "missing_in_a" {
		t.Errorf("expected rank-1 kind first, got %q", sorted[0].Kind)
	}
	if sorted[1].Kind != "missing_in_b" {
		t.Errorf("expected rank-2 kind second, got %q", sorted[1].Kind)
	}
	if sorted[2].Kind != "mismatch" || sorted[3].Kind != "mismatch" {
		t.Error("expected last two to be mismatch")
	}
	// within mismatch, keys should be sorted
	if sorted[2].Key > sorted[3].Key {
		t.Errorf("secondary sort by key failed: %q > %q", sorted[2].Key, sorted[3].Key)
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	results := sampleResults()
	originalFirst := results[0].Key
	sorter.Sort(results, sorter.SortByKey)
	if results[0].Key != originalFirst {
		t.Error("Sort mutated the input slice")
	}
}

func TestGroupByKind(t *testing.T) {
	results := sampleResults()
	groups := sorter.GroupByKind(results)

	if len(groups["mismatch"]) != 2 {
		t.Errorf("expected 2 mismatches, got %d", len(groups["mismatch"]))
	}
	if len(groups["missing_in_b"]) != 1 {
		t.Errorf("expected 1 missing_in_b, got %d", len(groups["missing_in_b"]))
	}
	if len(groups["missing_in_a"]) != 1 {
		t.Errorf("expected 1 missing_in_a, got %d", len(groups["missing_in_a"]))
	}
}
