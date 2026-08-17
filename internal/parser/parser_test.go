package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name        string
		format      Format
		input       string
		want        any
		err         error
		errContains string
	}{
		{
			name:   "valid json",
			format: FormatJSON,
			input:  `{"x":10.5}`,
			want:   map[string]any{"x": 10.5},
		},
		{
			name:   "valid yaml",
			format: FormatYAML,
			input:  `x: 10.5`,
			want:   map[string]any{"x": 10.5},
		},
		{
			name:        "error invalid json",
			format:      FormatJSON,
			input:       `{]`,
			err:         ErrFailedToParse,
			errContains: "json",
		},
		{
			name:        "error invalid yaml",
			format:      FormatYAML,
			input:       `}{`,
			err:         ErrFailedToParse,
			errContains: "yaml",
		},
		{
			name:   "error unsupported format",
			format: FormatUnknown,
			input:  `}{`,
			err:    ErrUnsupportedFormat,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse([]byte(tc.input), tc.format)

			require.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

			if len(tc.errContains) > 0 {
				assert.ErrorContains(t, err, tc.errContains)
			}
		})
	}
}

func TestDeduceFormatFromPath(t *testing.T) {
	cases := []struct {
		name        string
		path        string
		want        Format
		err         error
		errContains string
	}{
		{
			name: ".json ext",
			path: "file.json",
			want: FormatJSON,
		},
		{
			name: ".yaml ext",
			path: "file.yaml",
			want: FormatYAML,
		},
		{
			name: ".yml ext",
			path: "file.yml",
			want: FormatYAML,
		},
		{
			name: "nested path",
			path: filepath.Join("dir", "nested", "file.json"),
			want: FormatJSON,
		},
		{
			name:        "unsupported format",
			path:        "file.obj",
			want:        FormatUnknown,
			err:         ErrUnsupportedFormat,
			errContains: ".obj",
		},
		{
			name:        "no extension",
			path:        "file",
			want:        FormatUnknown,
			err:         ErrUnsupportedFormat,
			errContains: noExtensionMessage,
		},
		{
			name:        "empty path",
			path:        "",
			want:        FormatUnknown,
			err:         ErrUnsupportedFormat,
			errContains: noExtensionMessage,
		},
	}

	for _, tc := range cases {
		got, err := DeduceFormatFromPath(tc.path)

		require.ErrorIs(t, err, tc.err)
		assert.Equal(t, tc.want, got)

		if len(tc.errContains) > 0 {
			assert.ErrorContains(t, err, tc.errContains)
		}
	}
}
