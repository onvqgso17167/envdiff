// Package filter provides utilities to narrow down diff results
// produced by the diff package. Callers can restrict output to
// specific key prefixes, an explicit key allowlist, or a particular
// class of difference (missing vs. mismatched).
//
// Typical usage:
//
//	results := diff.Compare(envA, envB)
//	filtered := filter.Apply(results, filter.Options{
//		Prefix:      "APP_",
//		OnlyMissing: true,
//	})
package filter
