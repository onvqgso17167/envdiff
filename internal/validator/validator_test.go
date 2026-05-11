package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/validator"
)

func TestValidate_Nonempty(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"API_KEY": ""},
	}
	rules := []validator.Rule{{Key: "API_KEY", Kind: "nonempty"}}
	v := validator.Validate(envs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "API_KEY" {
		t.Errorf("unexpected key: %s", v[0].Key)
	}
}

func TestValidate_URL_Valid(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"BASE_URL": "https://example.com"},
	}
	rules := []validator.Rule{{Key: "BASE_URL", Kind: "url"}}
	v := validator.Validate(envs, rules)
	if len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}

func TestValidate_URL_Invalid(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"BASE_URL": "not-a-url"},
	}
	rules := []validator.Rule{{Key: "BASE_URL", Kind: "url"}}
	v := validator.Validate(envs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_Bool(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"DEBUG": "yes"},
	}
	rules := []validator.Rule{{Key: "DEBUG", Kind: "bool"}}
	v := validator.Validate(envs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_Prefix(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"SECRET": "sk_live_abc"},
	}
	rules := []validator.Rule{{Key: "SECRET", Kind: "prefix:sk_"}}
	v := validator.Validate(envs, rules)
	if len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}

func TestValidate_Prefix_Fail(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"SECRET": "pk_live_abc"},
	}
	rules := []validator.Rule{{Key: "SECRET", Kind: "prefix:sk_"}}
	v := validator.Validate(envs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidateResults_Mismatch(t *testing.T) {
	results := []diff.Result{
		{Key: "DEBUG", Kind: "mismatch", ValueA: "yes", ValueB: "true"},
	}
	rules := []validator.Rule{{Key: "DEBUG", Kind: "bool"}}
	v := validator.ValidateResults(results, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for 'yes', got %d", len(v))
	}
}

func TestValidate_KeyNotPresent(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"OTHER": "value"},
	}
	rules := []validator.Rule{{Key: "MISSING", Kind: "nonempty"}}
	v := validator.Validate(envs, rules)
	if len(v) != 0 {
		t.Errorf("expected no violations for absent key, got %d", len(v))
	}
}
