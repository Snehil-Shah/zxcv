// Package install implements the `zxcv install` command.
package install

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/installer"
)

// Command returns the `install` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "install",
		Usage:       "Install resolved tools or add new ones",
		Description: "Downloads & installs all resolved tools for the current directory, if no arguments are provided. If a specific tool and its version are provided, it installs only that and adds them to the local .tool-versions manifest. This behavior can be controlled to save it to the global manifest or not save at all using the --global and --no-save flags, respectively. You can specify the version as 'latest' (or with a prefix 'latest:<prefix>') to simply install and pin the latest version of the tool.",
		ArgsUsage:   "[<tool> {<version>|latest|latest:<prefix>}]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "global", Usage: "Install globally. Affects the global ~/.tool-versions file"},
			&cli.BoolFlag{Name: "no-save", Usage: "Do not modify any .tool-versions file"},
		},
		Action: run,
	}
}

// run is the action function for the `install` command.
func run(ctx context.Context, c *cli.Command) error {
	args := c.Args().Slice()
	inst := installer.New()
	switch len(args) {
	case 0:
		if c.Bool("global") || c.Bool("no-save") {
			return errors.New("--global and --no-save require <tool> <version> args")
		}
		return runFromManifest(ctx, inst)
	case 2:
		return runSpecific(ctx, inst, c, args[0], args[1])
	default:
		return errors.New("usage: zxcv install [options] [<tool> <version>]")
	}
}
