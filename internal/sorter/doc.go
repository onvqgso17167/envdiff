// Package sorter provides sorting and grouping utilities for collections
// of diff.Result values produced by the envdiff comparison engine.
//
// Supported sort orders:
//
//   - SortByKey      — alphabetical by environment variable name (default)
//   - SortByKind     — alphabetical by result kind (missing_in_a, missing_in_b, mismatch)
//   - SortBySeverity — by impact rank: missing_in_a → missing_in_b → mismatch
//
// GroupByKind partitions a flat result slice into a map for structured
// reporting or further processing.
package sorter
