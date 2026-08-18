package json

import (
	"encoding/json"

	"code/internal/diff"
)

type diffItemJSON struct {
	Key      string         `json:"key"`
	Type     string         `json:"type"`
	Value1   any            `json:"value1,omitempty"`
	Value2   any            `json:"value2,omitempty"`
	Children []diffItemJSON `json:"children,omitempty"`
}

// Format formats diff to json representation.
func Format(df []diff.Item) string {
	rootChildren := convertDiffToJSON(df)

	jsonDiff := diffItemJSON{
		Key:      "",
		Type:     "root",
		Children: rootChildren,
	}

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
		Key:  item.Key,
		Type: getHumanReadbleChange(item),
	}

	if item.Nested != nil {
		jsonDiff := convertDiffToJSON(*item.Nested)

		diffItem.Children = jsonDiff

		return diffItem
	}

	if item.Change == diff.ItemChangeAdd {
		diffItem.Value2 = item.NewValue

		return diffItem
	}

	if item.Change == diff.ItemChangeRemove {
		diffItem.Value1 = item.OldValue

		return diffItem
	}

	if item.Change == diff.ItemChangeReplace {
		diffItem.Value1 = item.OldValue
		diffItem.Value2 = item.NewValue

		return diffItem
	}

	diffItem.Value1 = item.OldValue

	return diffItem
}

func getHumanReadbleChange(item diff.Item) string {
	if item.Nested != nil {
		return "nested"
	}

	switch item.Change {
	case diff.ItemChangeAdd:
		return "added"
	case diff.ItemChangeRemove:
		return "deleted"
	case diff.ItemChangeReplace:
		return "changed"
	case diff.ItemChangeNone:
		return "unchanged"
	default:
		return "unchanged"
	}
}
