package diff

import (
	"fmt"
	"slices"
	"strings"
)

// ItemChange represents one of possible changes: add, remove, no change.
type ItemChange int

// ItemChange option
const (
	ItemChangeNone = iota
	ItemChangeAdd
	ItemChangeRemove
)

// Item is item of internal diff representation
type Item struct {
	key    string
	value  any
	change ItemChange
}

// ComputeDiff computes diff and returns internal representation
func ComputeDiff(a, b map[string]any) []Item {
	result := make([]Item, 0, len(a)+len(b))

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
			result = append(result, Item{
				key:    k,
				value:  v2,
				change: ItemChangeAdd,
			})

			continue
		}

		if !ok2 {
			result = append(result, Item{
				key:    k,
				value:  v1,
				change: ItemChangeRemove,
			})

			continue
		}

		if v1 != v2 {
			result = append(result, Item{
				key:    k,
				value:  v1,
				change: ItemChangeRemove,
			})

			result = append(result, Item{
				key:    k,
				value:  v2,
				change: ItemChangeAdd,
			})

			continue
		}

		result = append(result, Item{
			key:    k,
			value:  v1,
			change: ItemChangeNone,
		})
	}

	return result
}

// FormatDiff takes internal diff representation and output formatted string
func FormatDiff(diff []Item) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for _, item := range diff {
		sb.WriteString("  ")
		sb.WriteString(getChangeSymbol(item.change))
		sb.WriteString(" ")
		fmt.Fprintf(&sb, "%s: %v", item.key, item.value)
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}

func getChangeSymbol(change ItemChange) string {
	switch change {
	case ItemChangeAdd:
		return "+"
	case ItemChangeRemove:
		return "-"
	case ItemChangeNone:
		return " "
	default:
		return " "
	}
}
