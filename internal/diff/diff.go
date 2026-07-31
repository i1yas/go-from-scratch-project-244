package diff

import (
	"fmt"
	"slices"
	"strings"
)

// GetDiff returns string with json-like diff report
func GetDiff(a, b map[string]any) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	keys := make([]string, 0, len(a)+len(b))
	for k := range a {
		keys = append(keys, k)
	}

	for k := range b {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}

	slices.Sort(keys)

	for _, k := range keys {
		v1, ok1 := a[k]
		v2, ok2 := b[k]

		if !ok1 {
			sb.WriteString("  - ")
			sb.WriteString(fmt.Sprintf("%s: %v", k, v2))
			sb.WriteString("\n")

			continue
		}

		if !ok2 {
			sb.WriteString("  - ")
			sb.WriteString(fmt.Sprintf("%s: %v", k, v1))
			sb.WriteString("\n")

			continue
		}

		if v1 != v2 {
			sb.WriteString("  - ")
			sb.WriteString(fmt.Sprintf("%s: %v", k, v1))
			sb.WriteString("\n")

			sb.WriteString("  + ")
			sb.WriteString(fmt.Sprintf("%s: %v", k, v2))
			sb.WriteString("\n")

			continue
		}

		sb.WriteString("    ")
		sb.WriteString(fmt.Sprintf("%s: %v", k, v1))
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
