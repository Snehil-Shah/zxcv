// Package exec implements the `zxcv exec` command.
package exec

import (
	"context"
	"errors"
	"strings"

	"github.com/urfave/cli/v3"
)

// Command returns the `exec` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:            "exec",
		Usage:           "Run a specific binary",
		Description:     "Runs the binary under a specified tool and version. If no version is specified for the tool, it defaults to the resolved version for the current working directory.",
		ArgsUsage:       "<tool>[@<version>] <binary> [args...]",
		SkipFlagParsing: true,
		Action:          run,
	}
}

// run is the action function for the `exec` command.
func run(ctx context.Context, c *cli.Command) error {
	args := c.Args().Slice()
	if len(args) < 2 {
		return errors.New("usage: zxcv exec <tool>[@<version>] <binary> [args...]")
	}
	pluginAtVer, binary, rest := args[0], args[1], args[2:]
	if strings.Contains(pluginAtVer, "@") {
		return runExplicit(ctx, pluginAtVer, binary, rest)
	}
	return runResolved(ctx, pluginAtVer, binary, rest)
}
