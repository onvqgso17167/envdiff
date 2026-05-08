package exporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Kind: diff.KindMismatch, ValueA: "localhost", ValueB: "prod.db"},
		{Key: "API_KEY", Kind: diff.KindMissingInB, ValueA: "abc123", ValueB: ""},
	}
}

func TestExport_CSV(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(sampleResults(), "dev.env", "prod.env", exporter.FormatCSV, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,kind,dev.env,prod.env") {
		t.Errorf("expected CSV header, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "mismatch") {
		t.Errorf("expected kind 'mismatch' in output, got:\n%s", out)
	}
}

func TestExport_Markdown(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(sampleResults(), "dev.env", "prod.env", exporter.FormatMarkdown, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Kind |") {
		t.Errorf("expected Markdown header, got:\n%s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "missing_in_b") {
		t.Errorf("expected kind 'missing_in_b' in output, got:\n%s", out)
	}
}

func TestExport_EmptyResults(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export([]diff.Result{}, "a.env", "b.env", exporter.FormatCSV, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(sampleResults(), "a.env", "b.env", exporter.Format("xml"), &buf)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
