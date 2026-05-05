package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestReport_TextNoDifferences(t *testing.T) {
	var buf bytes.Buffer
	Report(nil, FormatText, &buf)
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestReport_TextMissingInA(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.MissingInA},
	}
	var buf bytes.Buffer
	Report(results, FormatText, &buf)
	out := buf.String()
	if !strings.Contains(out, "[MISSING IN A]") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestReport_TextMismatch(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Status: diff.Mismatched, ValueA: "8080", ValueB: "9090"},
	}
	var buf bytes.Buffer
	Report(results, FormatText, &buf)
	out := buf.String()
	if !strings.Contains(out, "[MISMATCH]") || !strings.Contains(out, "PORT") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("values not shown in output: %s", out)
	}
}

func TestReport_JSONNoDifferences(t *testing.T) {
	var buf bytes.Buffer
	Report(nil, FormatJSON, &buf)
	out := buf.String()
	if !strings.Contains(out, `"differences":[]`) {
		t.Errorf("expected empty JSON array, got: %s", out)
	}
}

func TestReport_JSONWithResults(t *testing.T) {
	results := []diff.Result{
		{Key: "SECRET", Status: diff.MissingInB, ValueA: "abc", ValueB: ""},
	}
	var buf bytes.Buffer
	Report(results, FormatJSON, &buf)
	out := buf.String()
	if !strings.Contains(out, `"key":"SECRET"`) {
		t.Errorf("expected key in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"status":"missing_in_b"`) {
		t.Errorf("expected status in JSON, got: %s", out)
	}
}
