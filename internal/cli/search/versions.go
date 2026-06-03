// Package search implements the `zxcv search` command.
package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// runVersionSearch lists versions of name matching version query.
func runVersionSearch(ctx context.Context, name, query string) error {
	prog := ui.NewProgress([]string{ui.Dim("updating registry"), ui.Dim("syncing plugin")})
	if err := registry.Refresh(ctx); err != nil {
		prog.Failed(0)
		return fmt.Errorf("registry: %w", err)
	}
	prog.Done(0)

	p, err := installer.EnsurePlugin(ctx, name, nil)
	if err != nil {
		prog.Failed(1)
		return fmt.Errorf("plugin: %w", err)
	}
	prog.Done(1)
	prog.Clear()

	if plugin.IsLatestString(query) {
		v, err := p.Latest(ctx, query)
		if err != nil {
			return fmt.Errorf("latest: %w", err)
		}
		ui.Write(v)
		return nil
	}

	versions, err := p.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("list-all: %w", err)
	}
	var matches []string
	for _, v := range versions {
		if strings.HasPrefix(v, query) {
			matches = append(matches, v)
		}
	}
	if len(matches) == 0 {
		return fmt.Errorf("no versions of %s matching %q", name, query)
	}
	for _, v := range matches {
		ui.Write(v)
	}
	return nil
}
