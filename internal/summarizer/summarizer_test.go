package summarizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summarizer"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "A", Kind: diff.Clean},
		{Key: "B", Kind: diff.MissingInB},
		{Key: "C", Kind: diff.MissingInA},
		{Key: "D", Kind: diff.Mismatch},
		{Key: "E", Kind: diff.Mismatch},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := summarizer.Summarize(makeResults(), "a.env", "b.env")

	if s.TotalKeys != 5 {
		t.Errorf("TotalKeys: want 5, got %d", s.TotalKeys)
	}
	if s.Clean != 1 {
		t.Errorf("Clean: want 1, got %d", s.Clean)
	}
	if s.MissingInA != 1 {
		t.Errorf("MissingInA: want 1, got %d", s.MissingInA)
	}
	if s.MissingInB != 1 {
		t.Errorf("MissingInB: want 1, got %d", s.MissingInB)
	}
	if s.Mismatched != 2 {
		t.Errorf("Mismatched: want 2, got %d", s.Mismatched)
	}
}

func TestSummarize_HasDiffs(t *testing.T) {
	s := summarizer.Summarize(makeResults(), "a.env", "b.env")
	if !s.HasDiffs {
		t.Error("HasDiffs: expected true")
	}
}

func TestSummarize_NoDiffs(t *testing.T) {
	results := []diff.Result{
		{Key: "X", Kind: diff.Clean},
		{Key: "Y", Kind: diff.Clean},
	}
	s := summarizer.Summarize(results, "a.env", "b.env")
	if s.HasDiffs {
		t.Error("HasDiffs: expected false for all-clean results")
	}
	if s.DiffCount() != 0 {
		t.Errorf("DiffCount: want 0, got %d", s.DiffCount())
	}
}

func TestSummarize_FileNames(t *testing.T) {
	s := summarizer.Summarize(nil, "staging.env", "prod.env")
	if s.FileA != "staging.env" || s.FileB != "prod.env" {
		t.Errorf("unexpected file names: %q %q", s.FileA, s.FileB)
	}
}

func TestSummarize_DiffCount(t *testing.T) {
	s := summarizer.Summarize(makeResults(), "a.env", "b.env")
	if got := s.DiffCount(); got != 4 {
		t.Errorf("DiffCount: want 4, got %d", got)
	}
}
