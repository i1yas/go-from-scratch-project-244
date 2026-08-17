package plain

import (
	"code/internal/diff"
	"fmt"
	"reflect"
	"strings"
)

func Format(df []diff.Item) string {
	result := format(df, []string{})
	result = strings.TrimSpace(result)

	return result
}

func format(df []diff.Item, parentPath []string) string {
	var sb strings.Builder

	for _, item := range df {
		path := append(parentPath, item.Key)
		pathFormatted := strings.Join(path, ".")

		if item.Nested != nil {
			sb.WriteString(format(*item.Nested, path))

			continue
		}

		if item.Change == diff.ItemChangeAdd {
			fmt.Fprintf(&sb, "Property '%s' was added with value: %s\n", pathFormatted, formatValue(item.NewValue))

			continue
		}

		if item.Change == diff.ItemChangeRemove {
			fmt.Fprintf(&sb, "Property '%s' was removed\n", pathFormatted)

			continue
		}

		if item.Change == diff.ItemChangeReplace {
			oldValue := formatValue(item.OldValue)
			newValue := formatValue(item.NewValue)
			fmt.Fprintf(&sb, "Property '%s' was updated. From %s to %s\n", pathFormatted, oldValue, newValue)

			continue
		}
	}

	return sb.String()
}

func formatValue(value any) string {
	if value == nil {
		return "null"
	}

	kind := reflect.ValueOf(value).Kind()

	if kind == reflect.Map || kind == reflect.Slice {
		return "[complex value]"
	}

	if kind == reflect.String {
		return fmt.Sprintf("'%s'", value)
	}

	return fmt.Sprint(value)
}
