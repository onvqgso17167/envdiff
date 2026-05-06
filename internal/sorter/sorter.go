// Package sorter provides utilities for sorting and grouping diff results
// by severity, key name, or environment source.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortBy defines the field to sort results by.
type SortBy string

const (
	SortByKey      SortBy = "key"
	SortByKind     SortBy = "kind"
	SortBySeverity SortBy = "severity"
)

// severityRank assigns a numeric rank to each diff kind for ordering.
func severityRank(kind string) int {
	switch kind {
	case "missing_in_a":
		return 1
	case "missing_in_b":
		return 2
	case "mismatch":
		return 3
	default:
		return 99
	}
}

// Sort returns a new slice of diff.Result sorted according to the given SortBy field.
// For equal primary keys, results are secondarily sorted by key name.
func Sort(results []diff.Result, by SortBy) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		switch by {
		case SortByKind:
			if out[i].Kind != out[j].Kind {
				return out[i].Kind < out[j].Kind
			}
			return out[i].Key < out[j].Key
		case SortBySeverity:
			ri := severityRank(out[i].Kind)
			rj := severityRank(out[j].Kind)
			if ri != rj {
				return ri < rj
			}
			return out[i].Key < out[j].Key
		default: // SortByKey
			return out[i].Key < out[j].Key
		}
	})

	return out
}

// GroupByKind partitions results into a map keyed by their Kind field.
func GroupByKind(results []diff.Result) map[string][]diff.Result {
	groups := make(map[string][]diff.Result)
	for _, r := range results {
		groups[r.Kind] = append(groups[r.Kind], r)
	}
	return groups
}
