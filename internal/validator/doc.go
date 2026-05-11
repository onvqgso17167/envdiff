// Package validator provides value-level validation for environment variables
// discovered during a diff. Rules can enforce constraints such as non-empty
// values, valid URLs, boolean flags, or required string prefixes.
//
// Rules are applied per-key across all loaded env files, and violations are
// returned as structured Violation values for use in reports or CLI output.
package validator
