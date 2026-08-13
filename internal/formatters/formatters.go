package formatters

import (
	"fmt"

	"code/internal/diff"
	"code/internal/formatters/stylish"
)

const (
	// FormatStylish is human readable json-like format
	FormatStylish = "stylish"
)

// FormatDiff takes internal diff representation and output formatted string
func FormatDiff(diff []diff.Item, format string) (string, error) {
	if format == FormatStylish {
		return stylish.Format(diff), nil
	}

	return "", fmt.Errorf("Unsupported format: %s", format)
}
