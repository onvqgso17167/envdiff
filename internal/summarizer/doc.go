// Package summarizer aggregates diff.Result slices into a human-readable
// Summary struct, providing counts of clean, missing, and mismatched keys
// across two compared environment files.
//
// Typical usage:
//
//	results := diff.Compare(envA, envB)
//	s := summarizer.Summarize(results, "staging.env", "prod.env")
//	if s.HasDiffs {
//		fmt.Printf("%d differences found\n", s.DiffCount())
//	}
package summarizer
