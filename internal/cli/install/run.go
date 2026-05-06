// Package install implements the `zxcv install` command.
package install

import (
	"context"
	"fmt"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// runInstall handles the shared UX and returns installation results.
func runInstall(ctx context.Context, inst *installer.Installer, targets []installer.Target) ([]installer.Result, error) {
	if needsRegistry(targets) {
		regProg := ui.NewProgress([]string{ui.Dim("updating registry")})
		if err := registry.Refresh(ctx); err != nil {
			regProg.Failed(0)
			return nil, fmt.Errorf("refresh: %w", err)
		}
		regProg.Done(0)
		regProg.Clear()
	}

	ui.Info(ui.Dim("installing:"))
	prog := ui.NewProgress(labelsFor(targets))
	inst.OnComplete = progressCallback(prog)
	results := inst.Install(ctx, targets)
	return results, nil
}

// needsRegistry reports whether any target lacks a Definition and would hit the registry.
func needsRegistry(targets []installer.Target) bool {
	for _, t := range targets {
		if t.Definition == nil {
			return true
		}
	}
	return false
}
