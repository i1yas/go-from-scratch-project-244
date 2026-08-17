package stylish

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"code/internal/diff"
)

const (
	changeSymbolAdd    = '+'
	changeSymbolRemove = '-'
	changeSymbolNone   = ' '
)

func Format(diff []diff.Item) string {
	return formatMapsDiff(diff, 0)
}

func formatMapsDiff(df []diff.Item, nesting int) string {
	var sb strings.Builder

	sb.WriteString(formatMapOpen())

	for _, item := range df {
		if item.Nested != nil {
			value := formatMapsDiff(*item.Nested, nesting+1)
			row := FormatMapKeyValueRow(changeSymbolNone, item.Key, value, nesting)
			sb.WriteString(row)

			continue
		}

		if item.Change == diff.ItemChangeAdd {
			value := formatValue(item.NewValue, nesting+1)
			row := FormatMapKeyValueRow(changeSymbolAdd, item.Key, value, nesting)
			sb.WriteString(row)

			continue
		}

		if item.Change == diff.ItemChangeRemove {
			value := formatValue(item.OldValue, nesting+1)
			row := FormatMapKeyValueRow(changeSymbolRemove, item.Key, value, nesting)
			sb.WriteString(row)

			continue
		}

		if item.Change == diff.ItemChangeReplace {
			oldValue := formatValue(item.OldValue, nesting+1)
			newValue := formatValue(item.NewValue, nesting+1)

			oldRow := FormatMapKeyValueRow(changeSymbolRemove, item.Key, oldValue, nesting)
			newRow := FormatMapKeyValueRow(changeSymbolAdd, item.Key, newValue, nesting)

			sb.WriteString(oldRow)
			sb.WriteString(newRow)

			continue
		}

		value := formatValue(item.OldValue, nesting+1)
		row := FormatMapKeyValueRow(changeSymbolNone, item.Key, value, nesting)
		sb.WriteString(row)
	}

	sb.WriteString(formatMapClose(nesting))

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

	sb.WriteString(formatMapOpen())

	keys := make([]string, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, k := range keys {
		v := value[k]

		value := formatValue(v, nesting+1)
		row := FormatMapKeyValueRow(changeSymbolNone, k, value, nesting)
		sb.WriteString(row)
	}

	sb.WriteString(formatMapClose(nesting))

	return sb.String()
}

func formatMapOpen() string {
	return "{\n"
}

func formatMapClose(nesting int) string {
	return strings.Repeat(" ", 4*nesting) + "}"
}

func FormatMapKeyValueRow(changeSymbol rune, key string, value string, nesting int) string {
	pad := strings.Repeat(" ", 4*nesting)
	return fmt.Sprintf("%s  %c %s: %s\n", pad, changeSymbol, key, value)
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
	sb.WriteString("]")

	return sb.String()
}
