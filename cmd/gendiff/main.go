package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"code/internal/diff"
	"code/internal/parser"
)

var (
	ErrExpectingTwoPaths    = errors.New("expecting 2 paths")
	ErrExpectingRegularFile = errors.New("expecting regular file")
)

func main() {
	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<path> <path>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Usage:       "output format",
				Aliases:     []string{"f"},
				DefaultText: "stylish",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				return fmt.Errorf("%w: got %d", ErrExpectingTwoPaths, cmd.Args().Len())
			}

			data1, err := readFileFromArgument(cmd.Args().Get(0))
			if err != nil {
				return err
			}

			data2, err := readFileFromArgument(cmd.Args().Get(1))
			if err != nil {
				return err
			}

			parsed1, err := parser.Parse(data1)
			if err != nil {
				return err
			}

			parsed2, err := parser.Parse(data2)
			if err != nil {
				return err
			}

			fmt.Println(diff.GetDiff(parsed1, parsed2))

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
