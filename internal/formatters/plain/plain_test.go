package plain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
)

func TestFormat_EmptyDiff(t *testing.T) {
	input := []diff.Item{}
	want := ""

	got := Format(input)

	require.Equal(t, want, got)
}

func TestFormatValue(t *testing.T) {
	cases := []struct {
		name  string
		input any
		want  string
	}{
		{name: "string", input: "hello", want: "'hello'"},
		{name: "empty string", input: "", want: "''"},
		{name: "integer", input: int(5), want: "5"},
		{name: "float", input: float32(10.5), want: "10.5"},
		{name: "boolean", input: true, want: "true"},
		{name: "nil", input: nil, want: "null"},
		{name: "slice", input: []any{1, 2, 3}, want: "[complex value]"},
		{name: "map", input: map[string]any{"x": 10}, want: "[complex value]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatValue(tc.input)

			require.Equal(t, tc.want, got)
		})
	}
}
