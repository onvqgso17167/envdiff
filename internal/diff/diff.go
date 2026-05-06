package diff

import "sort"

// Status describes the kind of difference found for a key.
type Status string

const (
	MissingInA Status = "missing_in_a"
	MissingInB Status = "missing_in_b"
	Mismatch    Status = "mismatch"
	Match       Status = "match"
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key     string `json:"key"`
	Status  Status `json:"status"`
	ValueA  string `json:"value_a,omitempty"`
	ValueB  string `json:"value_b,omitempty"`
}

// Compare compares two env maps and returns a sorted list of Results.
// Only non-matching entries are returned unless includeMatches is true.
func Compare(a, b map[string]string) []Result {
	keys := unionKeys(a, b)
	sort.Strings(keys)

	var results []Result
	for _, k := range keys {
		va, inA := a[k]
		vb, inB := b[k]

		switch {
		case !inA:
			results = append(results, Result{Key: k, Status: MissingInA, ValueB: vb})
		case !inB:
			results = append(results, Result{Key: k, Status: MissingInB, ValueA: va})
		case va != vb:
			results = append(results, Result{Key: k, Status: Mismatch, ValueA: va, ValueB: vb})
		}
	}
	return results
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
