// Package current implements the `zxcv current` command.
package current

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/resolver"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// Command returns the `current` subcommand definition.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "current",
		Usage:       "Show all resolved tools and their versions for the current directory",
		Description: "Shows all active tools (with version and source) for the current working directory by resolving all applicable .tool-versions files and global installations.",
		Action:      run,
	}
}

// run is the action function for the `current` command.
func run(_ context.Context, _ *cli.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	resolutions, err := resolver.Resolve(cwd)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if len(resolutions) == 0 {
		ui.Info("no tools resolved — no .tool-versions file found in this directory, any parent, or $HOME")
		return nil
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "TOOL\tVERSION\tSOURCE")
	for _, r := range resolutions {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.Version, r.Source)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	ui.Write(strings.TrimRight(buf.String(), "\n"))
	return nil
}
