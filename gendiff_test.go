package code

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/formatters"
)

func TestGenDiff(t *testing.T) {
	jsonFile1 := getFixturePath(t, "nested1.json")
	jsonFile2 := getFixturePath(t, "nested2.json")

	yamlFile1 := getFixturePath(t, "nested1.yaml")
	yamlFile2 := getFixturePath(t, "nested2.yaml")

	arrayJSONFile := getFixturePath(t, "array.json")
	invalidJSONFile := getFixturePath(t, "invalid.json")

	stylishOutput := loadFixture(t, "nested1_nested2_stylish.txt")
	plainOutput := loadFixture(t, "nested1_nested2_plain.txt")

	cases := []struct {
		name   string
		input1 string
		input2 string
		format string
		want   string
		err    error
	}{
		{
			name:   "json files, stylish output",
			input1: jsonFile1,
			input2: jsonFile2,
			format: formatters.FormatStylish,
			want:   stylishOutput,
		},
		{
			name:   "yaml files, stylish output",
			input1: yamlFile1,
			input2: yamlFile2,
			format: formatters.FormatStylish,
			want:   stylishOutput,
		},
		{
			name:   "json files, plain output",
			input1: jsonFile1,
			input2: jsonFile2,
			format: formatters.FormatPlain,
			want:   plainOutput,
		},
		{
			name:   "error failed to deduce format",
			input1: "file-with-no-ext-1",
			input2: jsonFile2,
			format: formatters.FormatStylish,
			want:   "",
			err:    ErrFailedToDeduceFormat,
		},
		{
			name:   "error failed to deduce format",
			input1: jsonFile1,
			input2: "file-dose-not-exist.json",
			format: formatters.FormatStylish,
			want:   "",
			err:    ErrFailedToGetFileInfo,
		},
		{
			name:   "error failed to parse file",
			input1: invalidJSONFile,
			input2: jsonFile2,
			format: formatters.FormatStylish,
			want:   "",
			err:    ErrFailedToParseFile,
		},
		{
			name:   "error expecting map",
			input1: arrayJSONFile,
			input2: jsonFile2,
			format: formatters.FormatStylish,
			want:   "",
			err:    ErrExpectedMap,
		},
		{
			name:   "error failed to format diff",
			input1: jsonFile1,
			input2: jsonFile2,
			format: "some-unknown-format",
			want:   "",
			err:    ErrFailedToFormatDiff,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GenDiff(tc.input1, tc.input2, tc.format)

			require.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("json files, json output", func(t *testing.T) {
		jsonFile1 := getFixturePath(t, "nested1.json")
		jsonFile2 := getFixturePath(t, "nested2.json")

		var parsedWant, parsedGot any

		want := loadFixture(t, "nested1_nested2_json.txt")

		err := json.Unmarshal([]byte(want), &parsedWant)
		require.NoError(t, err)

		got, err := GenDiff(jsonFile1, jsonFile2, formatters.FormatJSON)
		require.NoError(t, err)

		err = json.Unmarshal([]byte(got), &parsedGot)
		require.NoError(t, err)

		assert.Equal(t, parsedWant, parsedGot)
	})

	t.Run("error failed to build abs path", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "nested")

		err := os.Mkdir(dir, 0o755)
		require.NoError(t, err)

		err = os.Chdir(dir)
		require.NoError(t, err)

		err = os.Remove(dir)
		require.NoError(t, err)

		_, err = GenDiff("file1.json", "file2.json", formatters.FormatStylish)

		require.ErrorIs(t, err, ErrFailedToBuildAbsPath)
	})

	t.Run("error expecting regular file", func(t *testing.T) {
		dirWithExt := filepath.Join(t.TempDir(), "dir.json")

		err := os.Mkdir(dirWithExt, 0o644)
		require.NoError(t, err)

		_, err = GenDiff(dirWithExt, jsonFile2, formatters.FormatStylish)

		require.ErrorIs(t, err, ErrExpectingRegularFile)
	})

	t.Run("error failed to read file", func(t *testing.T) {
		skipWhenUserIsRoot(t)

		forbiddenFile := filepath.Join(t.TempDir(), "forbidden.json")

		err := os.WriteFile(forbiddenFile, []byte("{}"), 0o000)
		require.NoError(t, err)

		_, err = GenDiff(forbiddenFile, jsonFile2, formatters.FormatStylish)

		require.ErrorIs(t, err, ErrFailedToReadFile)
	})
}

func getFixturePath(t *testing.T, name string) string {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("testdata", "fixture", name))
	require.NoError(t, err)

	return path
}

func loadFixture(t *testing.T, name string) string {
	t.Helper()

	path := getFixturePath(t, name)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(data)
}

func skipWhenUserIsRoot(t *testing.T) {
	t.Helper()

	u, err := user.Current()
	require.NoError(t, err)

	isRoot := u.Uid == "0"

	if isRoot {
		t.Skip("skipping test under root")
	}
}
