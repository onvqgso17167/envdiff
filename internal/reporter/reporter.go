package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable or JSON diff report to the given writer.
func Report(results []diff.Result, format Format, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	switch format {
	case FormatJSON:
		writeJSON(results, w)
	default:
		writeText(results, w)
	}
}

func writeText(results []diff.Result, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	for _, r := range results {
		switch r.Status {
		case diff.MissingInA:
			fmt.Fprintf(w, "[MISSING IN A] %s\n", r.Key)
		case diff.MissingInB:
			fmt.Fprintf(w, "[MISSING IN B] %s\n", r.Key)
		case diff.Mismatched:
			fmt.Fprintf(w, "[MISMATCH]     %s: %q != %q\n", r.Key, r.ValueA, r.ValueB)
		}
	}
}

func writeJSON(results []diff.Result, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, `{"differences":[]}`)
		return
	}

	var sb strings.Builder
	sb.WriteString(`{"differences":[`)
	for i, r := range results {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(
			`{"key":%q,"status":%q,"value_a":%q,"value_b":%q}`,
			r.Key, r.Status, r.ValueA, r.ValueB,
		))
	}
	sb.WriteString("]}") 
	fmt.Fprintln(w, sb.String())
}
