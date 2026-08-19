package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/app"
	"code/internal/formatters"
)

var binaryPath = ""

func TestSmoke_DefaultFormat(t *testing.T) {
	file1Path := getFixturePath(t, "file1.json")
	file2Path := getFixturePath(t, "file2.json")
	output := loadFixture(t, "result_stylish.txt")
	want := output + "\n"

	cmd := exec.Command(binaryPath, file1Path, file2Path)

	stdout, stderr, err := runCommand(cmd)
	exitCode := cmd.ProcessState.ExitCode()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	require.Equal(t, want, stdout)
	require.Empty(t, stderr)
}

func TestSmoke_SpecifyFormat(t *testing.T) {
	file1Path := getFixturePath(t, "file1.json")
	file2Path := getFixturePath(t, "file2.json")
	output := loadFixture(t, "result_plain.txt")
	want := output + "\n"

	cmd := exec.Command(binaryPath, file1Path, file2Path, "-f", formatters.FormatPlain)

	stdout, stderr, err := runCommand(cmd)
	exitCode := cmd.ProcessState.ExitCode()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	require.Equal(t, want, stdout)
	require.Empty(t, stderr)
}

func TestSmoke_ErrorNotEnoughArguments(t *testing.T) {
	cmd := exec.Command(binaryPath)

	stdout, stderr, err := runCommand(cmd)
	exitCode := cmd.ProcessState.ExitCode()

	require.Error(t, err)
	require.Equal(t, 1, exitCode)
	require.Empty(t, stdout)
	require.Contains(t, stderr, app.ErrExpectingTwoArgs.Error())
}

func TestSmoke_ErrorTooMuchArguments(t *testing.T) {
	cmd := exec.Command(binaryPath, "file1.json", "file2.json", "file3.json")

	stdout, stderr, err := runCommand(cmd)
	exitCode := cmd.ProcessState.ExitCode()

	require.Error(t, err)
	require.Equal(t, 1, exitCode)
	require.Empty(t, stdout)
	require.Contains(t, stderr, app.ErrExpectingTwoArgs.Error())
}

func TestSmoke_ErrorUnsupportedFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--unknown-flag")

	stdout, stderr, err := runCommand(cmd)
	exitCode := cmd.ProcessState.ExitCode()

	require.Error(t, err)
	require.Equal(t, 1, exitCode)
	require.Contains(t, stdout, "USAGE:")
	require.Contains(t, stderr, "unknown-flag")
}

func TestMain(m *testing.M) {
	bin, err := buildTestingBinary()
	if err != nil {
		fmt.Printf("Failed to build binary: %s", err.Error())
		os.Exit(1)
	}

	defer func() {
		err := os.Remove(bin)
		if err != nil {
			fmt.Printf("Failed to cleanup binary: %s", err.Error())
			os.Exit(1)
		}
	}()

	binaryPath = bin

	m.Run()
}

func buildTestingBinary() (string, error) {
	bin := filepath.Join(os.TempDir(), "hexlet-diff_test_bin")
	cmd := exec.Command("go", "build", "-o", bin, ".")

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return bin, nil
}

func runCommand(cmd *exec.Cmd) (string, string, error) {
	var out, errOut bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()

	return out.String(), errOut.String(), err
}

func getFixturePath(t *testing.T, name string) string {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixture", name))
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
