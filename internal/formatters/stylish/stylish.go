package stylish

import (
	"code/internal/diff"
	"fmt"
	"reflect"
	"strings"
)

func Format(diff []diff.Item) string {
	return formatMapsDiff(diff, 0)
}

func formatMapsDiff(diff []diff.Item, nesting int) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for _, item := range diff {
		sb.WriteString(strings.Repeat(" ", 4*nesting))
		sb.WriteString("  ")
		sb.WriteString(getChangeSymbol(item.Change))
		sb.WriteString(" ")

		if item.Nested != nil {
			fmt.Fprintf(&sb, "%s: %s", item.Key, formatMapsDiff(*item.Nested, nesting+1))
		} else {
			fmt.Fprintf(&sb, "%s: %s", item.Key, formatValue(item.Value, nesting+1))
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

func getChangeSymbol(change diff.ItemChange) string {
	switch change {
	case diff.ItemChangeAdd:
		return "+"
	case diff.ItemChangeRemove:
		return "-"
	case diff.ItemChangeNone:
		return " "
	default:
		return " "
	}
}
