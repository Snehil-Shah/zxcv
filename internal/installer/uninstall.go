// Package installer implements the install and uninstall pipelines.
package installer

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Snehil-Shah/zxcv/internal/binary"
	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
)

// HasInstalls reports whether any tool is currently installed.
func HasInstalls() bool {
	_, err := os.Stat(config.InstallsDir())
	return err == nil
}

// UninstallAll wipes every installed tool in parallel, then removes leftover shim/plugin/download dirs.
func UninstallAll(ctx context.Context) error {
	tools, err := os.ReadDir(config.InstallsDir())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return wipeTopLevel()
		}
		return err
	}

	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	for _, t := range tools {
		if !t.IsDir() {
			continue
		}
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if err := UninstallTool(ctx, name); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}(t.Name())
	}
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return wipeTopLevel()
}

// wipeTopLevel removes the top-level data dirs that may sit empty after per-tool cleanup.
func wipeTopLevel() error {
	for _, dir := range []string{config.InstallsDir(), config.ShimsDir(), config.PluginsDir(), config.DownloadsDir()} {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("remove %s: %w", dir, err)
		}
	}
	return nil
}

// UninstallTool removes every installed version of name.
func UninstallTool(ctx context.Context, name string) error {
	if _, err := os.Stat(filepath.Join(config.InstallsDir(), name)); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%s is not installed", name)
		}
		return fmt.Errorf("stat %s: %w", name, err)
	}
	versions, err := VersionsFor(name)
	if err != nil {
		return err
	}
	for _, v := range versions {
		if err := UninstallVersion(ctx, name, v); err != nil {
			return err
		}
	}
	return wipeTool(ctx, name)
}

// wipeTool nukes every per-tool dir (install, plugin, download).
func wipeTool(ctx context.Context, name string) error {
	if err := plugin.New(name).Remove(ctx); err != nil {
		return err
	}
	for _, dir := range []string{
		filepath.Join(config.InstallsDir(), name),
		filepath.Join(config.DownloadsDir(), name),
	} {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("remove %s: %w", dir, err)
		}
	}
	// HACK: This is a double-op in most normal scenarios where the uninstallation of the last version would already prune the required. This is a safety callback for partial or broken installations. My mind is too clogged to think of a clean refactor for this.
	return pruneOrphanShims()
}

// UninstallVersion removes one installed version.
func UninstallVersion(ctx context.Context, name, version string) error {
	installPath := filepath.Join(config.InstallsDir(), name, version)
	if _, err := os.Stat(installPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%s %s is not installed", name, version)
		}
		return fmt.Errorf("stat %s: %w", installPath, err)
	}
	p := plugin.New(name)
	if err := preUninstall(ctx, p, version, installPath); err != nil {
		return err
	}
	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf("remove %s: %w", installPath, err)
	}
	toolDir := filepath.Join(config.InstallsDir(), name)
	entries, err := os.ReadDir(toolDir)
	if err != nil || len(entries) > 0 {
		return nil
	}
	_ = os.Remove(toolDir)
	if err := p.Remove(ctx); err != nil {
		return err
	}
	return pruneOrphanShims()
}

// preUninstall fires bin/uninstall (no-op if absent).
func preUninstall(ctx context.Context, p *plugin.Plugin, version, installPath string) error {
	env := plugin.AsdfEnv{Version: version, InstallPath: installPath}.Map()
	if _, _, err := p.RunCallback(ctx, "uninstall", nil, env); err != nil && !errors.Is(err, plugin.ErrCallbackNotImplemented) {
		return fmt.Errorf("uninstall callback: %w", err)
	}
	return nil
}

// VersionsFor returns the installed versions of name, skipping in-progress .partial dirs.
func VersionsFor(name string) ([]string, error) {
	toolDir := filepath.Join(config.InstallsDir(), name)
	entries, err := os.ReadDir(toolDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", toolDir, err)
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) == ".partial" {
			continue
		}
		out = append(out, e.Name())
	}
	return out, nil
}

// pruneOrphanShims removes shim files whose plugin no longer has any installed version.
func pruneOrphanShims() error {
	shimsDir := config.ShimsDir()
	entries, err := os.ReadDir(shimsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		p, err := binary.PluginForShim(e.Name())
		if err != nil {
			continue
		}
		toolEntries, err := os.ReadDir(filepath.Join(config.InstallsDir(), p.Name))
		if err == nil && len(toolEntries) > 0 {
			continue
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		_ = os.Remove(filepath.Join(shimsDir, e.Name()))
	}
	return nil
}
