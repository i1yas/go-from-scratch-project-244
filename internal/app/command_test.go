package app

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestCommand(t *testing.T) {
	testdata, err := getTestdataPath()
	require.NoError(t, err)

	file1 := filepath.Join(testdata, "nested1.json")
	file2 := filepath.Join(testdata, "nested2.json")

	stylishResultPath := filepath.Join(testdata, "nested1_nested2_stylish.txt")
	plainResultPath := filepath.Join(testdata, "nested1_nested2_plain.txt")

	stylishResult, err := os.ReadFile(stylishResultPath)
	require.NoError(t, err)

	plainResult, err := os.ReadFile(plainResultPath)
	require.NoError(t, err)

	cases := []struct {
		name           string
		args           []string
		output         string
		outputContains string
		err            error
	}{
		{
			name:   "default behavior",
			args:   []string{file1, file2},
			output: string(stylishResult) + "\n",
			err:    nil,
		},
		{
			name:   "plain format",
			args:   []string{file1, file2, "-f", "plain"},
			output: string(plainResult) + "\n",
			err:    nil,
		},
		{
			name: "error not enough params",
			args: []string{file1},
			err:  ErrExpectingTwoArgs,
		},
		{
			name: "error too much params",
			args: []string{file1},
			err:  ErrExpectingTwoArgs,
		},
		{
			name: "error failed to generate diff",
			args: []string{"does_not_exist", "does_not_exist"},
			err:  ErrFailedToGenerateDiff,
		},
		{
			name:           "error unknown flag",
			args:           []string{file1, file2, "--unknown-flag"},
			outputContains: "USAGE:",
			err:            errors.New("unknown-flag"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer

			cmd := BuildCommand()

			cmd.Writer = &outputBuf
			cmd.ExitErrHandler = func(context.Context, *cli.Command, error) {
				// NOTE: override this behavior with noop to prevent os.Exit call
			}

			const commandName = "test"

			args := append([]string{commandName}, tc.args...)

			err = cmd.Run(t.Context(), args)

			if tc.err == nil {
				require.NoError(t, err)
			} else {
				// NOTE: can't use ErrorIs beacause cli.Exit(err, code) does not wrap given error
				require.ErrorContains(t, err, tc.err.Error())
			}

			if len(tc.outputContains) > 0 {
				assert.Contains(t, outputBuf.String(), tc.outputContains)
			} else {
				assert.Equal(t, tc.output, outputBuf.String())
			}
		})
	}
}

func TestCommand_FailedToPrint(t *testing.T) {
	cmd := BuildCommand()

	wantErrMessage := "broken_writer"
	output := brokenWriter{err: errors.New(wantErrMessage)}

	testdata, err := getTestdataPath()
	require.NoError(t, err)

	file1 := filepath.Join(testdata, "nested1.json")
	file2 := filepath.Join(testdata, "nested2.json")

	cmd.Writer = &output
	cmd.ExitErrHandler = func(context.Context, *cli.Command, error) {
		// NOTE: override this behavior with noop to prevent os.Exit call
	}

	gotErr := cmd.Run(t.Context(), []string{"test", file1, file2})

	require.ErrorContains(t, gotErr, ErrFailedToPrintOutput.Error())
	require.ErrorContains(t, gotErr, wantErrMessage)
}

func getTestdataPath() (string, error) {
	return filepath.Abs(filepath.Join("..", "..", "testdata", "fixture"))
}

type brokenWriter struct {
	err error
}

func (w *brokenWriter) Write([]byte) (n int, err error) {
	return 0, w.err
}
