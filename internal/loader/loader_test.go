package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestLoad_Basic(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	ef, err := loader.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", ef.Vars["KEY"])
	}
	if ef.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", ef.Vars["FOO"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := loader.Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_PathIsStored(t *testing.T) {
	p := writeTempEnv(t, "A=1\n")
	ef, err := loader.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Path != p {
		t.Errorf("expected path %q, got %q", p, ef.Path)
	}
}

func TestLoadMultiple_Success(t *testing.T) {
	p1 := writeTempEnv(t, "X=1\n")
	p2 := writeTempEnv(t, "Y=2\n")
	files, err := loader.LoadMultiple([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Vars["X"] != "1" {
		t.Errorf("expected X=1")
	}
	if files[1].Vars["Y"] != "2" {
		t.Errorf("expected Y=2")
	}
}

func TestLoadMultiple_ErrorOnMissing(t *testing.T) {
	p1 := writeTempEnv(t, "X=1\n")
	_, err := loader.LoadMultiple([]string{p1, "/no/such/file"})
	if err == nil {
		t.Fatal("expected error for missing file in LoadMultiple")
	}
}
