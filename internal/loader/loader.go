package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its path and parsed key-value pairs.
type EnvFile struct {
	Path string
	Vars map[string]string
}

// Load reads and parses a .env file from the given path.
// Returns an error if the file does not exist or cannot be parsed.
func Load(path string) (*EnvFile, error) {
	clean := filepath.Clean(path)

	if _, err := os.Stat(clean); os.IsNotExist(err) {
		return nil, fmt.Errorf("env file not found: %s", clean)
	}

	vars, err := parser.ParseFile(clean)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", clean, err)
	}

	return &EnvFile{
		Path: clean,
		Vars: vars,
	}, nil
}

// LoadMultiple loads multiple .env files and returns them in order.
// All files are attempted; the first error encountered is returned.
func LoadMultiple(paths []string) ([]*EnvFile, error) {
	files := make([]*EnvFile, 0, len(paths))
	for _, p := range paths {
		ef, err := Load(p)
		if err != nil {
			return nil, err
		}
		files = append(files, ef)
	}
	return files, nil
}
