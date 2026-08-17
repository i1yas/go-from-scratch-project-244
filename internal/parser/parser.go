package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Format format acceptable by parser
type Format string

// Format option
const (
	FormatJSON    = "json"
	FormatYAML    = "yaml"
	FormatUnknown = "unknown"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrFailedToParse     = errors.New("failed to parse")
)

const noExtensionMessage = "no extension"

// Parse parses file content
func Parse(fileContent []byte, format Format) (any, error) {
	if format == FormatJSON {
		return parseJSON(fileContent)
	}

	if format == FormatYAML {
		return parseYAML(fileContent)
	}

	return nil, fmt.Errorf("%w '%s'", ErrUnsupportedFormat, format)
}

func parseJSON(fileContent []byte) (any, error) {
	var result any

	err := json.Unmarshal(fileContent, &result)
	if err != nil {
		return nil, fmt.Errorf("%w json: %w", ErrFailedToParse, err)
	}

	return result, nil
}

func parseYAML(fileContent []byte) (any, error) {
	var result any

	err := yaml.Unmarshal(fileContent, &result)
	if err != nil {
		return nil, fmt.Errorf("%w yaml: %w", ErrFailedToParse, err)
	}

	return result, nil
}

// DeduceFormatFromPath takes path with file extension and determines format from it.
func DeduceFormatFromPath(path string) (Format, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		return FormatJSON, nil
	case ".yaml", ".yml":
		return FormatYAML, nil
	default:
		message := ext
		if len(ext) == 0 {
			message = noExtensionMessage
		}

		return FormatUnknown, fmt.Errorf("%w: %s", ErrUnsupportedFormat, message)
	}
}
