// Package install implements the `zxcv install` command.
package install

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/resolver"
)

// runFromManifest installs every tool/version declared in the resolved `.tool-versions` chain.
func runFromManifest(ctx context.Context, inst *installer.Installer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	resolutions, err := resolver.Resolve(cwd)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}
	if len(resolutions) == 0 {
		return errors.New("no .tool-versions found in this directory or any parent")
	}

	targets := make([]installer.Target, len(resolutions))
	for i, r := range resolutions {
		targets[i] = installer.Target{
			Tool:       r.Tool,
			Definition: r.Definition,
		}
	}
	results, err := runInstall(ctx, inst, targets)
	if err != nil {
		return err
	}
	return finalize(results)
}
