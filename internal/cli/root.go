// Package cli builds the CLI entrypoint.
package cli

import (
	"github.com/urfave/cli/v3"

	// Commands:
	"github.com/Snehil-Shah/zxcv/internal/cli/current"
	"github.com/Snehil-Shah/zxcv/internal/cli/exec"
	"github.com/Snehil-Shah/zxcv/internal/cli/install"
	"github.com/Snehil-Shah/zxcv/internal/cli/list"
	"github.com/Snehil-Shah/zxcv/internal/cli/search"
	"github.com/Snehil-Shah/zxcv/internal/cli/uninstall"
	"github.com/Snehil-Shah/zxcv/internal/cli/which"
)

// New returns the top-level cli.Command.
func New(ver string) *cli.Command {
	return &cli.Command{
		Name:                  "zxcv",
		Version:               ver,
		Usage:                 "A minimalist, fast .tool-versions manager",
		Description:           "zxcv is a dead simple tool-version manager — a cleaner, friendlier alternative to asdf that plugs into its ecosystem with a highly minimalistic command surface and enjoyably better performance.",
		EnableShellCompletion: true,
		ConfigureShellCompletionCommand: func(c *cli.Command) {
			c.Hidden = false
		},
		Commands: []*cli.Command{
			current.Command(),
			search.Command(),
			install.Command(),
			exec.Command(),
			which.Command(),
			list.Command(),
			uninstall.Command(),
		},
	}
}
