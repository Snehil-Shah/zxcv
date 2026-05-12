// Package version implements the `zxcv version` command.
package version

import (
	"context"

	"github.com/urfave/cli/v3"
)

// Command returns the `version` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Usage:       "Output the zxcv version",
		Description: "Too scared to run the version command? I won't do anything more than showing you the zxcv version (trust me).",
		Action:      run,
	}
}

// run is the action function for the `version` command.
func run(_ context.Context, c *cli.Command) error {
	cli.ShowVersion(c.Root())
	return nil
}
