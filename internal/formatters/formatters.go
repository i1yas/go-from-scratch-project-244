package formatters

import (
	"fmt"

	"code/internal/diff"
	"code/internal/formatters/plain"
	"code/internal/formatters/stylish"
)

const (
	// FormatStylish is human readable json-like format
	FormatStylish = "stylish"
	// FormatPlain is flat format that show only changes
	FormatPlain = "plain"
)

// FormatDiff takes internal diff representation and output formatted string
func FormatDiff(diff []diff.Item, format string) (string, error) {
	if format == FormatStylish {
		return stylish.Format(diff), nil
	}

	if format == FormatPlain {
		return plain.Format(diff), nil
	}

	return "", fmt.Errorf("Unsupported format: %s", format)
}
