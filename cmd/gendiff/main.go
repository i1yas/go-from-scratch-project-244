package main

import (
	"code/internal/diff"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code"
)

var ErrExpectingTwoPaths = errors.New("expecting 2 paths")

func main() {
	defaultFormat := diff.FormatStylish

	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<path> <path>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Usage:       "output format",
				Aliases:     []string{"f"},
				Value:       defaultFormat,
				DefaultText: defaultFormat,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				return fmt.Errorf("%w: got %d", ErrExpectingTwoPaths, cmd.Args().Len())
			}

			format := cmd.String("format")

			diff, err := code.GenDiff(cmd.Args().Get(0), cmd.Args().Get(1), format)
			if err != nil {
				return err
			}

			fmt.Println(diff)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
