package validator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadRulesFile reads a simple rules file where each line is:
//
//	KEY=KIND
//
// Lines beginning with '#' and blank lines are ignored.
func LoadRulesFile(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("validator: open rules file: %w", err)
	}
	defer f.Close()

	var rules []Rule
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("validator: rules file line %d: invalid format %q", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		kind := strings.TrimSpace(parts[1])
		if key == "" || kind == "" {
			return nil, fmt.Errorf("validator: rules file line %d: key and kind must not be empty", lineNum)
		}
		rules = append(rules, Rule{Key: key, Kind: kind})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("validator: scan rules file: %w", err)
	}
	return rules, nil
}
