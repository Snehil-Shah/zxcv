// Package plugin manages asdf plugins.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Remove fires bin/pre-plugin-remove then deletes the plugin directory.
func (p *Plugin) Remove(ctx context.Context) error {
	if err := p.preRemove(ctx); err != nil {
		return err
	}
	if err := os.RemoveAll(p.Dir); err != nil {
		return fmt.Errorf("remove %s: %w", p.Dir, err)
	}
	return nil
}

// preRemove fires bin/pre-plugin-remove (no-op if absent).
func (p *Plugin) preRemove(ctx context.Context) error {
	env := AsdfEnv{PluginPath: p.Dir}.Map()
	if _, _, err := p.RunCallback(ctx, "pre-plugin-remove", nil, env); err != nil && !errors.Is(err, ErrCallbackNotImplemented) {
		return fmt.Errorf("pre-plugin-remove: %w", err)
	}
	return nil
}
