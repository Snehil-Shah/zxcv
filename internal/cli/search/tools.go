// Package search implements the `zxcv search` command.
package search

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// runNameSearch searches the registry for tool names matching query.
func runNameSearch(ctx context.Context, query string) error {
	prog := ui.NewProgress([]string{ui.Dim("updating registry")})
	if err := registry.Refresh(ctx); err != nil {
		prog.Failed(0)
		return fmt.Errorf("registry: %w", err)
	}
	prog.Done(0)
	prog.Clear()

	names, err := registry.ListAll()
	if err != nil {
		return fmt.Errorf("list registry: %w", err)
	}
	q := strings.ToLower(query)
	var matches []string
	for _, name := range names {
		if strings.Contains(strings.ToLower(name), q) {
			matches = append(matches, name)
		}
	}
	if len(matches) == 0 {
		return fmt.Errorf("no tools matching %q in registry — add custom tools via .tool-definitions", query)
	}
	sort.Strings(matches)
	for _, name := range matches {
		ui.Write(name)
	}
	return nil
}
