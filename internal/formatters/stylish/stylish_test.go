package stylish

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
)

func TestFormat_EmptyDiff(t *testing.T) {
	input := []diff.Item{}
	want := "{\n}"

	got := Format(input)

	require.Equal(t, want, got)
}

func TestFormatValue(t *testing.T) {
	sliceValue := loadJSONFixture(t, "slice.json")
	sliceFormatted := loadFixture(t, "slice_formatted.txt")

	mapValue := loadJSONFixture(t, "map.json")
	mapFormatted := loadFixture(t, "map_formatted.txt")

	cases := []struct {
		name  string
		input any
		want  string
	}{
		{name: "string", input: "hello", want: "hello"},
		{name: "empty string", input: "", want: ""},
		{name: "integer", input: int(5), want: "5"},
		{name: "float", input: float32(10.5), want: "10.5"},
		{name: "boolean", input: true, want: "true"},
		{name: "nil", input: nil, want: "null"},
		{name: "slice", input: sliceValue, want: sliceFormatted},
		{name: "map", input: mapValue, want: mapFormatted},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatValue(tc.input, 0)

			require.Equal(t, tc.want, got)
		})
	}
}

func loadFixture(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("testdata", "fixture", name)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(data)
}

func loadJSONFixture(t *testing.T, name string) any {
	t.Helper()

	fileContent := loadFixture(t, name)

	var result any

	err := json.Unmarshal([]byte(fileContent), &result)
	require.NoError(t, err)

	return result
}
