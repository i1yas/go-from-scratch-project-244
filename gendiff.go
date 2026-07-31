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
	data1, err := readFileFromArgument(path1)
	if err != nil {
		return "", err
	}

	data2, err := readFileFromArgument(path2)
	if err != nil {
		return "", err
	}

	parsed1, err := parser.Parse(data1)
	if err != nil {
		return "", err
	}

	parsed2, err := parser.Parse(data2)
	if err != nil {
		return "", err
	}

	diff := diff.GetDiff(parsed1, parsed2)

	return diff, nil
}

func readFileFromArgument(arg string) ([]byte, error) {
	path, err := filepath.Abs(arg)
	if err != nil {
		return []byte{}, err
	}

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return []byte{}, err
	}

	if !fileInfo.Mode().IsRegular() {
		return []byte{}, fmt.Errorf("bad argument '%s': %w", arg, ErrExpectingRegularFile)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
