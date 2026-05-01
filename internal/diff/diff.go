package diff

import (
	"fmt"
	"sort"
)

// ChangeType describes whether a value was added, removed, or changed.
type ChangeType string

const (
	ChangeAdded   ChangeType = "ADDED"
	ChangeRemoved ChangeType = "REMOVED"
	ChangeChanged ChangeType = "CHANGED"
)

// Change represents a single diff entry between two YAML values.
type Change struct {
	Type     ChangeType `json:"type"`
	Path     string     `json:"path"`
	OldValue any        `json:"old_value,omitempty"`
	NewValue any        `json:"new_value,omitempty"`
	Risk     string     `json:"risk"`
	Impact   []string   `json:"impact"`
}

// Compare performs a recursive diff of two Helm values maps.
func Compare(oldMap, newMap map[string]any) []Change {
	var changes []Change
	compareMaps("", oldMap, newMap, &changes)
	return changes
}

func compareAny(path string, old, new any, changes *[]Change) {
	oldMap, oldIsMap := asMap(old)
	newMap, newIsMap := asMap(new)
	oldList, oldIsList := asList(old)
	newList, newIsList := asList(new)

	switch {
	case oldIsMap && newIsMap:
		compareMaps(path, oldMap, newMap, changes)
	case oldIsList && newIsList:
		compareLists(path, oldList, newList, changes)
	default:
		if !equal(old, new) {
			*changes = append(*changes, Change{
				Type:     ChangeChanged,
				Path:     path,
				OldValue: old,
				NewValue: new,
			})
		}
	}
}

func compareMaps(path string, old, new map[string]any, changes *[]Change) {
	for _, key := range sortedUnion(old, new) {
		cp := joinPath(path, key)
		oldVal, inOld := old[key]
		newVal, inNew := new[key]

		switch {
		case inOld && !inNew:
			*changes = append(*changes, Change{
				Type:     ChangeRemoved,
				Path:     cp,
				OldValue: oldVal,
			})
		case !inOld && inNew:
			*changes = append(*changes, Change{
				Type:     ChangeAdded,
				Path:     cp,
				NewValue: newVal,
			})
		default:
			compareAny(cp, oldVal, newVal, changes)
		}
	}
}

func compareLists(path string, old, new []any, changes *[]Change) {
	min := len(old)
	if len(new) < min {
		min = len(new)
	}
	for i := 0; i < min; i++ {
		compareAny(listItemPath(path, i), old[i], new[i], changes)
	}
	for i := min; i < len(new); i++ {
		*changes = append(*changes, Change{
			Type:     ChangeAdded,
			Path:     listItemPath(path, i),
			NewValue: new[i],
		})
	}
	for i := min; i < len(old); i++ {
		*changes = append(*changes, Change{
			Type:     ChangeRemoved,
			Path:     listItemPath(path, i),
			OldValue: old[i],
		})
	}
}

func asMap(v any) (map[string]any, bool) {
	if v == nil {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

func asList(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	l, ok := v.([]any)
	return l, ok
}

func equal(a, b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func sortedUnion(a, b map[string]any) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
