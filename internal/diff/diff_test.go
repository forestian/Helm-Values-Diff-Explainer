package diff_test

import (
	"testing"

	"helm-values-diff-explainer/internal/diff"
)

func TestCompareAdded(t *testing.T) {
	old := map[string]any{}
	new := map[string]any{"key": "value"}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeAdded {
		t.Errorf("expected ADDED, got %s", changes[0].Type)
	}
	if changes[0].Path != "key" {
		t.Errorf("expected path 'key', got %s", changes[0].Path)
	}
	if changes[0].NewValue != "value" {
		t.Errorf("expected NewValue='value', got %v", changes[0].NewValue)
	}
}

func TestCompareRemoved(t *testing.T) {
	old := map[string]any{"key": "value"}
	new := map[string]any{}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeRemoved {
		t.Errorf("expected REMOVED, got %s", changes[0].Type)
	}
	if changes[0].OldValue != "value" {
		t.Errorf("expected OldValue='value', got %v", changes[0].OldValue)
	}
}

func TestCompareChangedScalar(t *testing.T) {
	old := map[string]any{"key": "old"}
	new := map[string]any{"key": "new"}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeChanged {
		t.Errorf("expected CHANGED, got %s", changes[0].Type)
	}
	if changes[0].OldValue != "old" {
		t.Errorf("expected OldValue='old', got %v", changes[0].OldValue)
	}
	if changes[0].NewValue != "new" {
		t.Errorf("expected NewValue='new', got %v", changes[0].NewValue)
	}
}

func TestCompareUnchanged(t *testing.T) {
	old := map[string]any{"key": "value", "num": 42}
	new := map[string]any{"key": "value", "num": 42}
	changes := diff.Compare(old, new)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d: %+v", len(changes), changes)
	}
}

func TestCompareNestedPath(t *testing.T) {
	old := map[string]any{"parent": map[string]any{"child": "old"}}
	new := map[string]any{"parent": map[string]any{"child": "new"}}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "parent.child" {
		t.Errorf("expected path 'parent.child', got %s", changes[0].Path)
	}
}

func TestCompareDeepNestedPath(t *testing.T) {
	old := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "old",
			},
		},
	}
	new := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "new",
			},
		},
	}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "a.b.c" {
		t.Errorf("expected path 'a.b.c', got %s", changes[0].Path)
	}
}

func TestCompareListIndexChanged(t *testing.T) {
	old := map[string]any{"list": []any{"a", "b", "c"}}
	new := map[string]any{"list": []any{"a", "x", "c"}}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "list[1]" {
		t.Errorf("expected path 'list[1]', got %s", changes[0].Path)
	}
}

func TestCompareListItemAdded(t *testing.T) {
	old := map[string]any{"list": []any{"a"}}
	new := map[string]any{"list": []any{"a", "b"}}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeAdded {
		t.Errorf("expected ADDED, got %s", changes[0].Type)
	}
	if changes[0].Path != "list[1]" {
		t.Errorf("expected path 'list[1]', got %s", changes[0].Path)
	}
}

func TestCompareListItemRemoved(t *testing.T) {
	old := map[string]any{"list": []any{"a", "b"}}
	new := map[string]any{"list": []any{"a"}}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeRemoved {
		t.Errorf("expected REMOVED, got %s", changes[0].Type)
	}
	if changes[0].Path != "list[1]" {
		t.Errorf("expected path 'list[1]', got %s", changes[0].Path)
	}
}

func TestCompareListOfMaps(t *testing.T) {
	old := map[string]any{
		"envs": []any{
			map[string]any{"name": "FOO", "value": "old"},
		},
	}
	new := map[string]any{
		"envs": []any{
			map[string]any{"name": "FOO", "value": "new"},
		},
	}
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "envs[0].value" {
		t.Errorf("expected path 'envs[0].value', got %s", changes[0].Path)
	}
}

func TestCompareMultipleChanges(t *testing.T) {
	old := map[string]any{"a": "1", "b": "2", "c": "3"}
	new := map[string]any{"a": "1", "b": "changed", "d": "added"}
	changes := diff.Compare(old, new)
	// b changed, c removed, d added = 3 changes
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d: %+v", len(changes), changes)
	}
}
