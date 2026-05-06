// Package search implements the `zxcv search` command.
package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// Command returns the `search` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "search",
		Usage:       "Search the upstream tool registry by tool name",
		Description: "Quickly finds the tools you want by name, or verifies if it doesn't exist before adding custom tools via .tool-definitions.",
		ArgsUsage:   "<query>",
		Action:      run,
	}
}

// run is the action function for the `search` command.
func run(ctx context.Context, c *cli.Command) error {
	query := c.Args().First()
	if query == "" {
		return errors.New("query required")
	}

	prog := ui.NewProgress([]string{ui.Dim("updating registry")})
	if err := registry.Refresh(ctx); err != nil {
		prog.Failed(0)
		return fmt.Errorf("registry: %w", err)
	}
	prog.Done(0)
	prog.Clear()

	matches, err := registry.Search(query)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no tools matching %q in registry — add custom tools via .tool-definitions", query)
	}
	for _, name := range matches {
		ui.Write(name)
	}
	return nil
}
