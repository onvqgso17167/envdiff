// Package exporter provides functionality to export diff results
// to various file formats such as CSV and Markdown.
package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for export.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Export writes diff results to w in the specified format.
func Export(results []diff.Result, fileA, fileB string, format Format, w io.Writer) error {
	switch format {
	case FormatCSV:
		return writeCSV(results, fileA, fileB, w)
	case FormatMarkdown:
		return writeMarkdown(results, fileA, fileB, w)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func writeCSV(results []diff.Result, fileA, fileB string, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "kind", fileA, fileB}); err != nil {
		return err
	}
	for _, r := range results {
		row := []string{r.Key, string(r.Kind), r.ValueA, r.ValueB}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(results []diff.Result, fileA, fileB string, w io.Writer) error {
	header := fmt.Sprintf("| Key | Kind | %s | %s |\n", fileA, fileB)
	sep := "|-----|------|" + strings.Repeat("-", len(fileA)+2) + "|" + strings.Repeat("-", len(fileB)+2) + "|\n"

	if _, err := fmt.Fprint(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, sep); err != nil {
		return err
	}
	for _, r := range results {
		line := fmt.Sprintf("| %s | %s | %s | %s |\n", r.Key, string(r.Kind), r.ValueA, r.ValueB)
		if _, err := fmt.Fprint(w, line); err != nil {
			return err
		}
	}
	return nil
}
