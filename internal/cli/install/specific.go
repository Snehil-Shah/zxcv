// Package install implements the `zxcv install` command.
package install

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/resolver"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

// runSpecific installs name@version, and optionally writes to `.tool-versions` (global or local).
func runSpecific(ctx context.Context, inst *installer.Installer, c *cli.Command, name, version string) error {
	targets := []installer.Target{{Tool: toolversions.Tool{Name: name, Version: version}}}
	results, err := runInstall(ctx, inst, targets)
	if err != nil {
		return err
	}
	if err := finalize(results); err != nil {
		return err
	}
	if c.Bool("no-save") {
		return nil
	}

	concrete := results[0].Target.Version
	path, err := writeTarget(c.Bool("global"))
	if err != nil {
		return err
	}
	if err := toolversions.WriteEntry(path, name, concrete); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// writeTarget returns the `.tool-versions` path to update, honoring `--global`.
func writeTarget(global bool) (string, error) {
	if global {
		return resolver.GlobalManifest()
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	return resolver.ApplicableManifest(cwd), nil
}
