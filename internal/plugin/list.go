// Package plugin manages asdf plugins.
package plugin

import (
	"context"
	"strings"
)

// ListAll returns all available versions reported by the plugin.
func (p *Plugin) ListAll(ctx context.Context) ([]string, error) {
	stdout, _, err := p.RunCallback(ctx, "list-all", nil, nil)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(stdout)), nil
}
