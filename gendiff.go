package code

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"code/internal/diff"
	"code/internal/parser"
)

var ErrExpectingRegularFile = errors.New("expecting regular file")

// GenDiff returns strings with json-like diff of two files
func GenDiff(path1 string, path2 string) (string, error) {
	parsed1, err := parseFileFromArgument(path1)
	if err != nil {
		return "", err
	}

	parsed2, err := parseFileFromArgument(path2)
	if err != nil {
		return "", err
	}

	diffRaw := diff.ComputeDiff(parsed1, parsed2)
	diffFormatted := diff.FormatDiff(diffRaw)

	return diffFormatted, nil
}

func parseFileFromArgument(arg string) (map[string]any, error) {
	var zero map[string]any

	path, err := filepath.Abs(arg)
	if err != nil {
		return zero, err
	}

	format, err := parser.DeduceFormatFromPath(path)
	if err != nil {
		return zero, err
	}

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return zero, err
	}

	if !fileInfo.Mode().IsRegular() {
		return zero, fmt.Errorf("bad argument '%s': %w", arg, ErrExpectingRegularFile)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}

	parsed, err := parser.Parse(data, format)
	if err != nil {
		return zero, err
	}

	return parsed, nil
}
