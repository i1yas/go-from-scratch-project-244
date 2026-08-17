package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"code"
	"code/internal/formatters"
)

var ErrExpectingTwoPaths = errors.New("expecting 2 paths")

// BuildCommand builds configured instance of cli.Command
func BuildCommand() *cli.Command {
	defaultFormat := formatters.FormatStylish

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
		Action: genDiffAction,
	}

	return cmd
}

func genDiffAction(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 2 {
		message := fmt.Sprintf("%s: got %d", ErrExpectingTwoPaths, cmd.Args().Len())
		return cli.Exit(message, 1)
	}

	format := cmd.String("format")

	diff, err := code.GenDiff(cmd.Args().Get(0), cmd.Args().Get(1), format)
	if err != nil {
		return cli.Exit(err, 1)
	}

	_, err = fmt.Fprintln(cmd.Root().Writer, diff)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}
