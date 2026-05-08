// Package summarizer provides aggregated statistics over diff results.
package summarizer

import "github.com/user/envdiff/internal/diff"

// Summary holds aggregated counts and metadata about a diff run.
type Summary struct {
	TotalKeys    int
	MissingInA   int
	MissingInB   int
	Mismatched   int
	Clean        int
	HasDiffs     bool
	FileA        string
	FileB        string
}

// Summarize computes a Summary from a slice of diff.Result and the two file paths.
func Summarize(results []diff.Result, fileA, fileB string) Summary {
	s := Summary{
		FileA: fileA,
		FileB: fileB,
	}

	for _, r := range results {
		s.TotalKeys++
		switch r.Kind {
		case diff.MissingInA:
			s.MissingInA++
		case diff.MissingInB:
			s.MissingInB++
		case diff.Mismatch:
			s.Mismatched++
		case diff.Clean:
			s.Clean++
		}
	}

	s.HasDiffs = s.MissingInA > 0 || s.MissingInB > 0 || s.Mismatched > 0
	return s
}

// DiffCount returns the total number of non-clean results.
func (s Summary) DiffCount() int {
	return s.MissingInA + s.MissingInB + s.Mismatched
}
