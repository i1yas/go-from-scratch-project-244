package diff

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	FormatStylish = "stylish"
)

// FormatDiff takes internal diff representation and output formatted string
func FormatDiff(diff []Item, format string) (string, error) {
	if format == FormatStylish {
		return formatStylish(diff), nil
	}

	return "", fmt.Errorf("Unsupported format: %s", format)
}

func formatStylish(diff []Item) string {
	return formatMapsDiff(diff, 0)
}

func formatMapsDiff(diff []Item, nesting int) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for _, item := range diff {
		sb.WriteString(strings.Repeat(" ", 4*nesting))
		sb.WriteString("  ")
		sb.WriteString(getChangeSymbol(item.change))
		sb.WriteString(" ")

		if item.nested != nil {
			fmt.Fprintf(&sb, "%s: %s", item.key, formatMapsDiff(*item.nested, nesting+1))
		} else {
			fmt.Fprintf(&sb, "%s: %s", item.key, formatValue(item.value, nesting+1))
		}

		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat(" ", 4*nesting))
	sb.WriteString("}")

	return sb.String()
}

func formatValue(value any, nesting int) string {
	kind := reflect.ValueOf(value).Kind()

	if kind == reflect.Map {
		m := value.(map[string]any)
		return formatMap(m, nesting)
	}

	if kind == reflect.Slice {
		s := value.([]any)
		return formatSlice(s, nesting)
	}

	if value == nil {
		return "null"
	}

	return fmt.Sprint(value)
}

func formatMap(value map[string]any, nesting int) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for k, v := range value {
		sb.WriteString(strings.Repeat(" ", 4*(nesting+1)))

		fmt.Fprintf(&sb, "%s: %s", k, formatValue(v, nesting+1))

		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat(" ", 4*nesting))
	sb.WriteString("}")

	return sb.String()
}

func formatSlice(value []any, nesting int) string {
	var sb strings.Builder

	sb.WriteString("[\n")

	for _, item := range value {
		sb.WriteString(strings.Repeat(" ", 4*(nesting+1)))

		sb.WriteString(formatValue(item, nesting+1))

		sb.WriteString(",\n")
	}

	sb.WriteString(strings.Repeat(" ", 4*nesting))
	sb.WriteString("]\n")

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
