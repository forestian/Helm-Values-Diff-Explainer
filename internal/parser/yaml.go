package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile reads a YAML file and returns its contents as a map.
func ParseFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q: %w", path, err)
	}
	return Parse(data)
}

// Parse decodes YAML bytes into a map. Returns an empty map for empty input.
func Parse(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	if result == nil {
		result = map[string]any{}
	}
	return result, nil
}
