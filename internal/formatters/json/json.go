package json

import (
	"encoding/json"

	"code/internal/diff"
)

type diffItemJSON map[string]any

// Format formats diff to json representation.
func Format(df []diff.Item) string {
	jsonDiff := convertDiffToJSON(df)

	data, _ := json.Marshal(jsonDiff)

	return string(data)
}

func convertDiffToJSON(df []diff.Item) []diffItemJSON {
	result := make([]diffItemJSON, len(df))

	for i, item := range df {
		result[i] = convertDiffItemToJSON(item)
	}

	return result
}

func convertDiffItemToJSON(item diff.Item) diffItemJSON {
	diffItem := diffItemJSON{
		"key":    item.Key,
		"change": getHumanReadbleChange(item.Change),
	}

	if item.Nested != nil {
		jsonDiff := convertDiffToJSON(*item.Nested)

		diffItem["nested"] = jsonDiff

		return diffItem
	}

	if item.Change == diff.ItemChangeAdd {
		diffItem["value"] = item.NewValue

		return diffItem
	}

	if item.Change == diff.ItemChangeRemove {
		diffItem["value"] = item.OldValue

		return diffItem
	}

	if item.Change == diff.ItemChangeReplace {
		diffItem["old_value"] = item.OldValue
		diffItem["new_value"] = item.NewValue

		return diffItem
	}

	diffItem["value"] = item.OldValue

	return diffItem
}

func getHumanReadbleChange(change diff.ItemChange) string {
	switch change {
	case diff.ItemChangeAdd:
		return "add"
	case diff.ItemChangeRemove:
		return "remove"
	case diff.ItemChangeReplace:
		return "replace"
	case diff.ItemChangeNone:
		return "none"
	default:
		return "none"
	}
}
