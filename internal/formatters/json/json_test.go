package json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
)

const emptyDiffOutput = `
{
	"key": "",
	"type": "root"
}
`

func TestFormat_EmptyDiff(t *testing.T) {
	input := []diff.Item{}
	want := strings.TrimSpace(emptyDiffOutput)

	got := Format(input)

	var parsedWant, parsedGot any

	err := json.Unmarshal([]byte(got), &parsedGot)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(want), &parsedWant)
	require.NoError(t, err)

	require.Equal(t, parsedWant, parsedGot)
}
