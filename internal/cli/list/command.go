// Package list implements the `zxcv list` command.
package list

import (
	"context"

	"github.com/urfave/cli/v3"
)

// Command returns the `list` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List installed tools and their versions",
		Description: "Lists all the installed tools and their versions. Pass a tool name to see installed versions for a specific tool.",
		ArgsUsage:   "[<tool>]",
		Action:      run,
	}
}

// run is the action function for the `list` command.
func run(_ context.Context, c *cli.Command) error {
	if name := c.Args().First(); name != "" {
		return listOne(name)
	}
	return listAll()
}
