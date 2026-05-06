// Package search implements the `zxcv search` command.
package search

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
)

// Command returns the `search` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "search",
		Usage:       "Search the upstream tool registry by tool name, or a tool's available versions",
		Description: "Quickly finds the tools or a tool's available versions you want by name. You can use this to discover if a tool exists and by what name, or finding the right version to pin it to.",
		ArgsUsage:   "<tool-query> | <tool> {<version-query>|latest|latest:<prefix>}",
		Action:      run,
	}
}

// run is the action function for the `search` command.
func run(ctx context.Context, c *cli.Command) error {
	args := c.Args().Slice()
	switch len(args) {
	case 1:
		return runNameSearch(ctx, args[0])
	case 2:
		return runVersionSearch(ctx, args[0], args[1])
	default:
		return errors.New("usage: zxcv search <tool-query> | <tool> {<version-query>|latest|latest:<prefix>}")
	}
}
