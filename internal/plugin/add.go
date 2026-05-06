// Package plugin manages asdf plugins.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// Add clones url into the plugin directory and fires bin/post-plugin-add.
func (p *Plugin) Add(ctx context.Context, url string) error {
	if err := os.MkdirAll(filepath.Dir(p.Dir), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(p.Dir), err)
	}
	if _, err := git.PlainCloneContext(ctx, p.Dir, false, &git.CloneOptions{URL: url, Depth: 1}); err != nil {
		return fmt.Errorf("clone %s: %w", url, err)
	}
	return p.postAdd(ctx)
}

// Link points the plugin directory at a local path via symlink and fires bin/post-plugin-add.
func (p *Plugin) Link(ctx context.Context, path string) error {
	if err := validatePluginDir(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p.Dir), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(p.Dir), err)
	}
	if err := os.Symlink(path, p.Dir); err != nil {
		return fmt.Errorf("symlink %s: %w", p.Dir, err)
	}
	return p.postAdd(ctx)
}

// postAdd fires bin/post-plugin-add (no-op if absent).
func (p *Plugin) postAdd(ctx context.Context) error {
	env := AsdfEnv{PluginPath: p.Dir}.Map()
	if _, _, err := p.RunCallback(ctx, "post-plugin-add", nil, env); err != nil && !errors.Is(err, ErrCallbackNotImplemented) {
		return fmt.Errorf("post-plugin-add: %w", err)
	}
	return nil
}

// ValidateExistence checks if the plugin exists on disk.
func (p *Plugin) ValidateExistence() error {
	return validatePluginDir(p.Dir)
}

// validatePluginDir checks that path looks like an asdf plugin.
func validatePluginDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s: not a directory", path)
	}
	install := filepath.Join(path, "bin", "install")
	if _, err := os.Stat(install); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%s: missing bin/install", path)
		}
		return fmt.Errorf("stat %s: %w", install, err)
	}
	return nil
}
