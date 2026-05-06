package filter

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds filtering configuration.
type Options struct {
	OnlyMissing  bool
	OnlyMismatch bool
	Prefix       string
	Keys         []string
}

// Apply filters a slice of diff.Result according to the given Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	var out []diff.Result
	for _, r := range results {
		if !matchesKind(r, opts) {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, opts.Prefix) {
			continue
		}
		if len(opts.Keys) > 0 && !containsKey(opts.Keys, r.Key) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func matchesKind(r diff.Result, opts Options) bool {
	if !opts.OnlyMissing && !opts.OnlyMismatch {
		return true
	}
	if opts.OnlyMissing && (r.Status == diff.MissingInA || r.Status == diff.MissingInB) {
		return true
	}
	if opts.OnlyMismatch && r.Status == diff.Mismatch {
		return true
	}
	return false
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
