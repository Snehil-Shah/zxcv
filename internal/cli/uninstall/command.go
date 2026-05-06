// Package uninstall implements the `zxcv uninstall` command.
package uninstall

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// Command returns the `uninstall` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "uninstall",
		Usage:       "Remove installed tools",
		Description: "Uninstalls tools previously installed by zxcv. With no arguments, removes all tools. With a tool name, removes all versions of that tool. With a tool name and version, removes just that version.",
		ArgsUsage:   "[<tool> [<version>]]",
		Action:      run,
	}
}

// run is the action function for the `uninstall` command.
func run(ctx context.Context, c *cli.Command) error {
	args := c.Args().Slice()
	switch len(args) {
	case 0:
		if !installer.HasInstalls() {
			ui.Info("nothing to uninstall")
			return nil
		}
		return withProgress("uninstalling all tools", func() error { return installer.UninstallAll(ctx) })
	case 1:
		return withProgress("uninstalling "+args[0], func() error {
			return installer.UninstallTool(ctx, args[0])
		})
	case 2:
		return withProgress("uninstalling "+args[0]+"@"+args[1], func() error {
			return installer.UninstallVersion(ctx, args[0], args[1])
		})
	default:
		return errors.New("usage: zxcv uninstall [<tool> [<version>]]")
	}
}

// withProgress runs fn under a single-line progress spinner.
func withProgress(label string, fn func() error) error {
	prog := ui.NewProgress([]string{label})
	if err := fn(); err != nil {
		prog.Failed(0)
		return err
	}
	prog.Done(0)
	return nil
}
