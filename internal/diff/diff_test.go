package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestCompare_Clean(t *testing.T) {
	a := parser.EnvMap{"KEY": "value", "PORT": "8080"}
	b := parser.EnvMap{"KEY": "value", "PORT": "8080"}

	result := Compare(a, b)
	if !result.IsClean() {
		t.Errorf("expected clean result, got %+v", result)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := parser.EnvMap{"KEY": "value", "EXTRA": "only_in_a"}
	b := parser.EnvMap{"KEY": "value"}

	result := Compare(a, b)
	if len(result.MissingInB) != 1 || result.MissingInB[0] != "EXTRA" {
		t.Errorf("expected MissingInB=[EXTRA], got %v", result.MissingInB)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := parser.EnvMap{"KEY": "value"}
	b := parser.EnvMap{"KEY": "value", "EXTRA": "only_in_b"}

	result := Compare(a, b)
	if len(result.MissingInA) != 1 || result.MissingInA[0] != "EXTRA" {
		t.Errorf("expected MissingInA=[EXTRA], got %v", result.MissingInA)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	a := parser.EnvMap{"DB_HOST": "localhost", "PORT": "8080"}
	b := parser.EnvMap{"DB_HOST": "prod.db.example.com", "PORT": "8080"}

	result := Compare(a, b)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatched key, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "DB_HOST" {
		t.Errorf("expected mismatched key DB_HOST, got %q", m.Key)
	}
	if m.ValueA != "localhost" || m.ValueB != "prod.db.example.com" {
		t.Errorf("unexpected values: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	a := parser.EnvMap{"Z_KEY": "v", "A_KEY": "v", "M_KEY": "v"}
	b := parser.EnvMap{}

	result := Compare(a, b)
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, key := range result.MissingInB {
		if key != expected[i] {
			t.Errorf("expected sorted key %q at index %d, got %q", expected[i], i, key)
		}
	}
}

func TestResult_IsClean_False(t *testing.T) {
	r := Result{MissingInA: []string{"KEY"}}
	if r.IsClean() {
		t.Error("expected IsClean to return false")
	}
}
