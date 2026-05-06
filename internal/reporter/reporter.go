// Package reporter formats and writes diff results to an io.Writer.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

// Options controls how results are rendered.
type Options struct {
	Format  string // "text" or "json"
	SortBy  sorter.SortBy
	Grouped bool // group by kind in text output
}

// Report writes diff results to w using the provided options.
// Returns an error if writing fails.
func Report(w io.Writer, results []diff.Result, opts Options) error {
	if opts.SortBy == "" {
		opts.SortBy = sorter.SortByKey
	}

	sorted := sorter.Sort(results, opts.SortBy)

	switch opts.Format {
	case "json":
		return writeJSON(w, sorted)
	default:
		return writeText(w, sorted, opts.Grouped)
	}
}

func writeText(w io.Writer, results []diff.Result, grouped bool) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	if grouped {
		groups := sorter.GroupByKind(results)
		for _, kind := range []string{"missing_in_a", "missing_in_b", "mismatch"} {
			group, ok := groups[kind]
			if !ok {
				continue
			}
			if _, err := fmt.Fprintf(w, "[%s]\n", kind); err != nil {
				return err
			}
			for _, r := range group {
				if err := writeTextLine(w, r); err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, r := range results {
		if err := writeTextLine(w, r); err != nil {
			return err
		}
	}
	return nil
}

func writeTextLine(w io.Writer, r diff.Result) error {
	switch r.Kind {
	case "missing_in_a":
		_, err := fmt.Fprintf(w, "MISSING_IN_A  %s\n", r.Key)
		return err
	case "missing_in_b":
		_, err := fmt.Fprintf(w, "MISSING_IN_B  %s\n", r.Key)
		return err
	default:
		_, err := fmt.Fprintf(w, "MISMATCH      %s (%q vs %q)\n", r.Key, r.ValueA, r.ValueB)
		return err
	}
}

func writeJSON(w io.Writer, results []diff.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
