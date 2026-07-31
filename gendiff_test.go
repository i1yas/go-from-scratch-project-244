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
	input1, err := getFixturePath(t, "test1.json")
	require.NoError(t, err)

	input2, err := getFixturePath(t, "test2.json")
	require.NoError(t, err)

	wantData, err := loadFixture(t, "test1_test2_result.txt")
	require.NoError(t, err)

	want := strings.TrimSpace(string(wantData))

	got, err := GenDiff(input1, input2)
	require.NoError(t, err)

	assert.Equal(t, string(want), got)
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
