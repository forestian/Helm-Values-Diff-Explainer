package parser_test

import (
	"testing"

	"helm-values-diff-explainer/internal/parser"
)

func TestParseValidYAML(t *testing.T) {
	data := []byte("key: value\nnested:\n  inner: 42\n")
	result, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result["key"])
	}
	nested, ok := result["nested"].(map[string]any)
	if !ok {
		t.Fatal("expected nested to be a map")
	}
	if nested["inner"] != 42 {
		t.Errorf("expected nested.inner=42, got %v", nested["inner"])
	}
}

func TestParseInvalidYAML(t *testing.T) {
	data := []byte("key: [\ninvalid yaml here")
	_, err := parser.Parse(data)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParseEmptyYAML(t *testing.T) {
	data := []byte("")
	result, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseNullYAML(t *testing.T) {
	data := []byte("~")
	result, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map for null YAML, got %v", result)
	}
}
