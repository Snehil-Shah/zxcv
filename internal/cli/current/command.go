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
		Usage:       "Show resolved tools and their versions for the current directory",
		Description: "Shows all active tools (with version and source) for the current working directory by resolving all applicable .tool-versions files and global installations. If a specific tool is provided, prints only its resolution.",
		ArgsUsage:   "[<tool>]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Usage: "output version(s) only"},
		},
		Action: run,
	}
}

// run is the action function for the `current` command.
func run(_ context.Context, c *cli.Command) error {
	tool := c.Args().First()
	versionOnly := c.Bool("version")

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	resolutions, err := resolver.Resolve(cwd)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if tool != "" {
		var match *resolver.Resolution
		for i := range resolutions {
			if resolutions[i].Name == tool {
				match = &resolutions[i]
				break
			}
		}
		if match == nil {
			return fmt.Errorf("no resolved version for tool %q", tool)
		}
		if versionOnly {
			// No table needed.
			ui.Write(match.Version)
			return nil
		}
		resolutions = []resolver.Resolution{*match}
	}

	if len(resolutions) == 0 {
		ui.Info("no tools resolved — no .tool-versions file found in this directory, any parent, or $HOME")
		return nil
	}

	showTool := tool == ""
	showSource := !versionOnly

	var headers, cols []string
	if showTool {
		headers = append(headers, "TOOL")
	}
	headers = append(headers, "VERSION")
	if showSource {
		headers = append(headers, "SOURCE")
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, r := range resolutions {
		cols = cols[:0]
		if showTool {
			cols = append(cols, r.Name)
		}
		cols = append(cols, r.Version)
		if showSource {
			cols = append(cols, r.Source)
		}
		_, _ = fmt.Fprintln(w, strings.Join(cols, "\t"))
	}
	if err := w.Flush(); err != nil {
		return err
	}
	ui.Write(strings.TrimRight(buf.String(), "\n"))
	return nil
}
