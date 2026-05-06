// Package which implements the `zxcv which` command.
package which

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/binary"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// Command returns the `which` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:      "which",
		Usage:     "Output the resolved path of a tool's binary",
		ArgsUsage: "<name>",
		Action:    run,
	}
}

// run is the action function for the `which` command.
func run(ctx context.Context, c *cli.Command) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("usage: zxcv which <name>")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	p, err := binary.PluginForShim(name)
	if err != nil {
		return err
	}
	r, err := binary.Find(ctx, p.Name, name, cwd)
	if err != nil {
		return err
	}
	ui.Write(r.BinPath)
	return nil
}
