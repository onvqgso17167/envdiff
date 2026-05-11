// Package validator provides value-level validation rules for env diff results.
// It checks that values conform to expected patterns (e.g. non-empty, URL, boolean).
package validator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Rule defines a named validation rule applied to a key's value.
type Rule struct {
	Key     string
	Kind    string // "nonempty", "url", "bool", "prefix:<val>"
}

// Violation describes a single validation failure.
type Violation struct {
	File    string
	Key     string
	Value   string
	Message string
}

// Validate applies the given rules against the loaded env maps and returns any violations.
func Validate(envs map[string]map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for file, env := range envs {
		for _, rule := range rules {
			val, exists := env[rule.Key]
			if !exists {
				continue
			}
			if msg := applyRule(val, rule.Kind); msg != "" {
				violations = append(violations, Violation{
					File:    file,
					Key:     rule.Key,
					Value:   val,
					Message: msg,
				})
			}
		}
	}
	return violations
}

// ValidateResults checks diff results against rules, flagging mismatched or missing keys.
func ValidateResults(results []diff.Result, rules []Rule) []Violation {
	var violations []Violation
	for _, r := range results {
		for _, rule := range rules {
			if r.Key != rule.Key {
				continue
			}
			for _, val := range []string{r.ValueA, r.ValueB} {
				if val == "" {
					continue
				}
				if msg := applyRule(val, rule.Kind); msg != "" {
					violations = append(violations, Violation{
						Key:     r.Key,
						Value:   val,
						Message: msg,
					})
				}
			}
		}
	}
	return violations
}

func applyRule(val, kind string) string {
	switch {
	case kind == "nonempty":
		if strings.TrimSpace(val) == "" {
			return "value must not be empty"
		}
	case kind == "url":
		if _, err := url.ParseRequestURI(val); err != nil {
			return fmt.Sprintf("value %q is not a valid URL", val)
		}
	case kind == "bool":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Sprintf("value %q is not a boolean", val)
		}
	case strings.HasPrefix(kind, "prefix:"):
		expected := strings.TrimPrefix(kind, "prefix:")
		if !strings.HasPrefix(val, expected) {
			return fmt.Sprintf("value %q must start with %q", val, expected)
		}
	}
	return ""
}
