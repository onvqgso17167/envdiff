package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
	"github.com/user/envdiff/internal/sorter"
)

func TestReport_TextNoDifferences(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.Report(&buf, nil, reporter.Options{Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestReport_TextMissingInA(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Kind: "missing_in_a"},
	}
	var buf bytes.Buffer
	err := reporter.Report(&buf, results, reporter.Options{Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "MISSING_IN_A") || !strings.Contains(buf.String(), "FOO") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReport_TextMismatch(t *testing.T) {
	results := []diff.Result{
		{Key: "BAR", Kind: "mismatch", ValueA: "old", ValueB: "new"},
	}
	var buf bytes.Buffer
	err := reporter.Report(&buf, results, reporter.Options{Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISMATCH") || !strings.Contains(out, "BAR") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestReport_JSONNoDifferences(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.Report(&buf, []diff.Result{}, reporter.Options{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array, got: %q", buf.String())
	}
}

func TestReport_JSONWithResults(t *testing.T) {
	results := []diff.Result{
		{Key: "X", Kind: "mismatch", ValueA: "1", ValueB: "2"},
	}
	var buf bytes.Buffer
	err := reporter.Report(&buf, results, reporter.Options{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"Key\"") {
		t.Errorf("expected JSON key field, got: %q", buf.String())
	}
}

func TestReport_GroupedText(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Kind: "mismatch", ValueA: "1", ValueB: "2"},
		{Key: "B", Kind: "missing_in_b"},
		{Key: "C", Kind: "missing_in_a"},
	}
	var buf bytes.Buffer
	err := reporter.Report(&buf, results, reporter.Options{
		Format:  "text",
		Grouped: true,
		SortBy:  sorter.SortBySeverity,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[missing_in_a]") {
		t.Errorf("expected group header [missing_in_a], got: %q", out)
	}
	if !strings.Contains(out, "[mismatch]") {
		t.Errorf("expected group header [mismatch], got: %q", out)
	}
}
