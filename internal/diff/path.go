package diff

import "fmt"

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

func listItemPath(parent string, index int) string {
	return fmt.Sprintf("%s[%d]", parent, index)
}
