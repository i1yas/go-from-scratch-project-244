package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Foramt format acceptable by parser
type Foramt string

// Foramt option
const (
	ForamtJSON    = "json"
	ForamtYAML    = "yaml"
	ForamtUnknown = "unknown"
)

var ErrUnsupportedFormat = errors.New("unsupported format")

// Parse parses file content
func Parse(fileContent []byte, format Foramt) (map[string]any, error) {
	if format == ForamtJSON {
		return parseJSON(fileContent)
	}

	if format == ForamtYAML {
		return parseYAML(fileContent)
	}

	return map[string]any{}, fmt.Errorf("unsupported format '%s'", format)
}

func parseJSON(fileContent []byte) (map[string]any, error) {
	var result map[string]any

	err := json.Unmarshal(fileContent, &result)
	if err != nil {
		return result, fmt.Errorf("failed to parse json: %w", err)
	}

	return result, nil
}

func parseYAML(fileContent []byte) (map[string]any, error) {
	var result map[string]any

	err := yaml.Unmarshal(fileContent, &result)
	if err != nil {
		return result, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return result, nil
}

// DeduceFormatFromPath takes path with file extension and determines format from it.
func DeduceFormatFromPath(path string) (Foramt, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		return ForamtJSON, nil
	case ".yaml", ".yml":
		return ForamtYAML, nil
	default:
		return ForamtUnknown, fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}
}
