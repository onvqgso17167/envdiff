package diff

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the outcome of comparing two EnvMaps.
type Result struct {
	// MissingInB contains keys present in A but absent in B.
	MissingInB []string
	// MissingInA contains keys present in B but absent in A.
	MissingInA []string
	// Mismatched contains keys present in both but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key    string
	ValueA string
	ValueB string
}

// IsClean returns true when there are no differences between the two env maps.
func (r Result) IsClean() bool {
	return len(r.MissingInA) == 0 &&
		len(r.MissingInB) == 0 &&
		len(r.Mismatched) == 0
}

// Compare analyses two EnvMaps and returns a Result describing their differences.
func Compare(a, b parser.EnvMap) Result {
	result := Result{}

	for key, valA := range a {
		valB, ok := b[key]
		if !ok {
			result.MissingInB = append(result.MissingInB, key)
			continue
		}
		if valA != valB {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				ValueA: valA,
				ValueB: valB,
			})
		}
	}

	for key := range b {
		if _, ok := a[key]; !ok {
			result.MissingInA = append(result.MissingInA, key)
		}
	}

	sort.Strings(result.MissingInA)
	sort.Strings(result.MissingInB)
	sort.Slice(result.Mismatched, func(i, j int) bool {
		return result.Mismatched[i].Key < result.Mismatched[j].Key
	})

	return result
}
