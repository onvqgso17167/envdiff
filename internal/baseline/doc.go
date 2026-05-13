// Package baseline manages reference snapshots of envdiff results.
//
// A baseline snapshot captures the diff results at a known point in time,
// allowing subsequent runs to detect newly introduced issues (regressions)
// rather than reporting all pre-existing differences as failures.
//
// Typical usage:
//
//	// Save a baseline after a clean review:
//	_ = baseline.Save(".envdiff-baseline.json", "dev.env", "prod.env", results)
//
//	// On the next run, load the baseline and find only new problems:
//	snap, _ := baseline.Load(".envdiff-baseline.json")
//	newProblems := baseline.NewIssues(snap, currentResults)
package baseline
