package main

import (
	"context"
	"errors"
	"os"

	"code/internal/app"
)

var ErrExpectingTwoPaths = errors.New("expecting 2 paths")

func main() {
	cmd := app.BuildCommand()

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(1)
	}
}
