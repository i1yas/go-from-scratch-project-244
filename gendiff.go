package code

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parser"
)

var (
	ErrFailedToBuildAbsPath = errors.New("failed to build absolute path")
	ErrFailedToDeduceFormat = errors.New("failed to deduce format")
	ErrFailedToGetFileInfo  = errors.New("failed to get file info")
	ErrExpectingRegularFile = errors.New("expecting regular file")
	ErrFailedToReadFile     = errors.New("failed to read file")
	ErrFailedToParseFile    = errors.New("failed to parse file")
	ErrExpectedMap          = errors.New("expected map")
	ErrFailedToFormatDiff   = errors.New("failed to format diff")
)

// GenDiff returns strings with json-like diff of two files
func GenDiff(path1, path2, format string) (string, error) {
	parsed1, err := parseFileFromArgument(path1)
	if err != nil {
		return "", err
	}

	parsed2, err := parseFileFromArgument(path2)
	if err != nil {
		return "", err
	}

	diffRaw := diff.ComputeDiff(parsed1, parsed2)

	diffFormatted, err := formatters.FormatDiff(diffRaw, format)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToFormatDiff, err)
	}

	return diffFormatted, nil
}

func parseFileFromArgument(arg string) (map[string]any, error) {
	var zero map[string]any

	path, err := filepath.Abs(arg)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrFailedToBuildAbsPath, err)
	}

	format, err := parser.DeduceFormatFromPath(path)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrFailedToDeduceFormat, err)
	}

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrFailedToGetFileInfo, err)
	}

	if !fileInfo.Mode().IsRegular() {
		return zero, fmt.Errorf("bad argument '%s': %w", arg, ErrExpectingRegularFile)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
	}

	parsed, err := parser.Parse(data, format)
	if err != nil {
		return zero, fmt.Errorf("%w '%s': %w", ErrFailedToParseFile, path, err)
	}

	parsedMap, ok := parsed.(map[string]any)
	if !ok {
		kind := reflect.ValueOf(parsed).Kind()
		return zero, fmt.Errorf("%w in '%s', got %s", ErrExpectedMap, path, kind.String())
	}

	return parsedMap, nil
}
