package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	t.Run("basic json", func(t *testing.T) {
		input1, err := getFixturePath(t, "flat1.json")
		require.NoError(t, err)

		input2, err := getFixturePath(t, "flat2.json")
		require.NoError(t, err)

		wantData, err := loadFixture(t, "flat1_flat2_result.txt")
		require.NoError(t, err)

		want := strings.TrimSpace(string(wantData))

		got, err := GenDiff(input1, input2)
		require.NoError(t, err)

		assert.Equal(t, string(want), got)
	})

	t.Run("basic yaml", func(t *testing.T) {
		input1, err := getFixturePath(t, "flat1.yaml")
		require.NoError(t, err)

		input2, err := getFixturePath(t, "flat2.yaml")
		require.NoError(t, err)

		wantData, err := loadFixture(t, "flat1_flat2_result.txt")
		require.NoError(t, err)

		want := strings.TrimSpace(string(wantData))

		got, err := GenDiff(input1, input2)
		require.NoError(t, err)

		assert.Equal(t, string(want), got)
	})

	t.Run("file does not exist", func(t *testing.T) {
		badInput, err := getFixturePath(t, "this_file_does_not_exist")
		require.NoError(t, err)

		input, err := getFixturePath(t, "flat1.json")
		require.NoError(t, err)

		_, err = GenDiff(badInput, input)
		require.Error(t, err)

		_, err = GenDiff(input, badInput)
		require.Error(t, err)
	})

	t.Run("path is not file", func(t *testing.T) {
		badInput, err := getFixturePath(t, ".")
		require.NoError(t, err)

		input, err := getFixturePath(t, "flat1.json")
		require.NoError(t, err)

		_, err = GenDiff(badInput, input)
		require.Error(t, err)

		_, err = GenDiff(input, badInput)
		require.Error(t, err)
	})
}

func getFixturePath(t *testing.T, name string) (string, error) {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("testdata", "fixture", name))
	if err != nil {
		return "", fmt.Errorf("failed to fixture path for '%s': %w", name, err)
	}

	return path, nil
}

func loadFixture(t *testing.T, name string) ([]byte, error) {
	t.Helper()

	path, err := getFixturePath(t, name)
	if err != nil {
		return []byte{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read fixture '%s', %w", name, err)
	}

	return data, nil
}
