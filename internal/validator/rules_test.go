package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func writeTempRules(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envrules")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempRules: %v", err)
	}
	return p
}

func TestLoadRulesFile_Basic(t *testing.T) {
	p := writeTempRules(t, "API_KEY=nonempty\nBASE_URL=url\nDEBUG=bool\n")
	rules, err := validator.LoadRulesFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Key != "API_KEY" || rules[0].Kind != "nonempty" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[2].Key != "DEBUG" || rules[2].Kind != "bool" {
		t.Errorf("unexpected rule[2]: %+v", rules[2])
	}
}

func TestLoadRulesFile_CommentsAndBlanks(t *testing.T) {
	p := writeTempRules(t, "# this is a comment\n\nSECRET=prefix:sk_\n")
	rules, err := validator.LoadRulesFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Kind != "prefix:sk_" {
		t.Errorf("unexpected kind: %s", rules[0].Kind)
	}
}

func TestLoadRulesFile_InvalidLine(t *testing.T) {
	p := writeTempRules(t, "BADLINE\n")
	_, err := validator.LoadRulesFile(p)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestLoadRulesFile_NotFound(t *testing.T) {
	_, err := validator.LoadRulesFile("/nonexistent/.envrules")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
