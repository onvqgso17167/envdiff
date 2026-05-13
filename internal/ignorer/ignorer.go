// Package ignorer provides functionality to load and apply ignore rules
// for envdiff comparisons. Keys matching ignore patterns are excluded from
// diff results, similar to a .gitignore for environment variables.
package ignorer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// IgnoreList holds a set of key patterns to exclude from diff results.
type IgnoreList struct {
	patterns []string
}

// LoadFile reads an ignore file where each non-blank, non-comment line
// is treated as an exact key name or prefix pattern ending with '*'.
func LoadFile(path string) (*IgnoreList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ignorer: open %q: %w", path, err)
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ignorer: read %q: %w", path, err)
	}
	return &IgnoreList{patterns: patterns}, nil
}

// New creates an IgnoreList from a slice of patterns directly.
func New(patterns []string) *IgnoreList {
	cp := make([]string, len(patterns))
	copy(cp, patterns)
	return &IgnoreList{patterns: cp}
}

// Matches reports whether the given key matches any pattern in the list.
// Patterns ending with '*' are treated as prefix matches.
func (il *IgnoreList) Matches(key string) bool {
	for _, p := range il.patterns {
		if strings.HasSuffix(p, "*") {
			if strings.HasPrefix(key, strings.TrimSuffix(p, "*")) {
				return true
			}
		} else if p == key {
			return true
		}
	}
	return false
}

// Apply filters out any diff.Result whose key is matched by the IgnoreList.
func (il *IgnoreList) Apply(results []diff.Result) []diff.Result {
	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if !il.Matches(r.Key) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
