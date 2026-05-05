package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func buildBinary(t *testing.T) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = "."
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, b)
	}
	return out
}

func TestCLI_NoDifferences(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "a.env", "KEY=val\n")
	b := writeTempEnv(t, "b.env", "KEY=val\n")
	cmd := exec.Command(bin, a, b)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "No differences") {
		t.Errorf("expected 'No differences' in output, got: %s", out)
	}
}

func TestCLI_WithDifferences_ExitCode2(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "a.env", "KEY=val\n")
	b := writeTempEnv(t, "b.env", "OTHER=val\n")
	cmd := exec.Command(bin, a, b)
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got: %v", err)
	}
}

func TestCLI_MissingFile_ExitCode1(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "/no/such/file", "/also/missing")
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got: %v", err)
	}
}
