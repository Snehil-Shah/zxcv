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
	"strings"
	"sync"

	"github.com/Snehil-Shah/zxcv/internal/binary"
	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
)

// Install runs the install pipeline for each target in parallel.
func (i *Installer) Install(ctx context.Context, targets []Target) []Result {
	if len(targets) == 0 {
		return nil
	}
	sem := make(chan struct{}, runtime.NumCPU()) // NOTE: added safety net as our install pipeline exponentially shells out bash processes.
	results := make([]Result, len(targets))
	var wg sync.WaitGroup
	for idx, target := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			res := installTarget(ctx, target)
			results[idx] = res
			if i.OnComplete != nil {
				i.progressMu.Lock()
				i.OnComplete(idx, res)
				i.progressMu.Unlock()
			}
		}()
	}
	wg.Wait()
	return results
}

// installTarget runs the install pipeline for a single target.
func installTarget(ctx context.Context, target Target) Result {
	res := Result{Target: target}

	def, err := sourceFor(target)
	if err != nil {
		res.Err = err
		return res
	}
	p := plugin.New(def.Name)
	if err := syncPlugin(ctx, p, def); err != nil {
		res.Err = fmt.Errorf("plugin: %w", err)
		return res
	}

	if isLatest(target.Version) {
		concrete, err := p.Latest(ctx, latestPrefix(target.Version))
		if err != nil {
			res.Err = fmt.Errorf("latest: %w", err)
			return res
		}
		target.Version = concrete
		res.Target = target
	}

	installPath := filepath.Join(config.InstallsDir(), target.Name, target.Version)
	if _, err := os.Stat(installPath); err == nil {
		// Existing installation found...
		if err := binary.Generate(ctx, p, target.Version, installPath); err != nil {
			res.Err = fmt.Errorf("shim: %w", err)
		}
		return res
	}

	// Let's start the process. Starting with download:
	// Pre-clean: wipe stale state from any prior interrupted run.
	downloadPath := filepath.Join(config.DownloadsDir(), target.Name, target.Version)
	_ = os.RemoveAll(downloadPath)
	if err := os.MkdirAll(downloadPath, 0o755); err != nil {
		res.Err = fmt.Errorf("mkdir downloads: %w", err)
		return res
	}
	defer func() {
		_ = os.RemoveAll(downloadPath)
		_ = os.Remove(filepath.Dir(downloadPath)) // also remove the per-tool parent if no other versions left it populated.
	}()

	partialPath := installPath + ".partial"
	_ = os.RemoveAll(partialPath)
	if err := os.MkdirAll(partialPath, 0o755); err != nil {
		res.Err = fmt.Errorf("mkdir installs: %w", err)
		return res
	}
	defer func() { _ = os.RemoveAll(partialPath) }() // cleans up on any failure before the rename below... no-op afterwards since the path is gone.

	env := callbackEnv(target, partialPath, downloadPath)

	out, errout, err := p.RunCallback(ctx, "download", nil, env)
	res.Stdout = append(res.Stdout, out...)
	res.Stderr = append(res.Stderr, errout...)
	if err != nil && !errors.Is(err, plugin.ErrCallbackNotImplemented) {
		res.Err = fmt.Errorf("download: %w", err)
		return res
	}

	// Install:
	out, errout, err = p.RunCallback(ctx, "install", nil, env)
	res.Stdout = append(res.Stdout, out...)
	res.Stderr = append(res.Stderr, errout...)
	if err != nil {
		res.Err = fmt.Errorf("install: %w", err)
		return res
	}

	if err := os.Rename(partialPath, installPath); err != nil {
		res.Err = fmt.Errorf("rename %s: %w", partialPath, err)
		return res
	}

	// Shim generation:
	if err := binary.Generate(ctx, p, target.Version, installPath); err != nil {
		res.Err = fmt.Errorf("shim: %w", err)
		return res
	}
	return res
}

// TODO: Maybe sourcing and plugin management should be part of plugin package? idk. currently the plugin package is just a dumb wrapper around its disk existence.

// sourceFor returns the plugin source for target, falling back to the registry.
func sourceFor(target Target) (*tooldefinitions.Definition, error) {
	if target.Definition != nil {
		return target.Definition, nil
	}
	url, found, err := registry.Lookup(target.Name)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("%s not in registry — add it to .tool-definitions or check the name by using the search command", target.Name)
	}
	return &tooldefinitions.Definition{Name: target.Name, URL: url}, nil
}

// syncPlugin reconciles the on-disk plugin with def, swapping the source when the
// type changed (URL vs path) or the symlink target moved.
// TODO: URL-to-URL changes aren't detected.
func syncPlugin(ctx context.Context, p *plugin.Plugin, def *tooldefinitions.Definition) error {
	info, err := os.Lstat(p.Dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return addPlugin(ctx, p, def)
		}
		return fmt.Errorf("stat %s: %w", p.Dir, err)
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0
	needsSwap := false
	switch {
	case def.Path != "" && !isSymlink:
		needsSwap = true
	case def.Path != "" && isSymlink:
		target, err := os.Readlink(p.Dir)
		if err != nil {
			return fmt.Errorf("readlink %s: %w", p.Dir, err)
		}
		needsSwap = target != def.Path
	case def.URL != "" && isSymlink:
		needsSwap = true
	case def.URL != "" && !isSymlink:
		// Detect a partial clone from an interrupted previous run.
		if err := p.ValidateExistence(); err != nil {
			needsSwap = true
		}
	}

	if !needsSwap {
		return nil
	}
	if err := p.Remove(ctx); err != nil {
		return err
	}
	return addPlugin(ctx, p, def)
}

// addPlugin clones or symlinks the plugin based on def.
func addPlugin(ctx context.Context, p *plugin.Plugin, def *tooldefinitions.Definition) error {
	if def.URL != "" {
		return p.Add(ctx, def.URL)
	}
	return p.Link(ctx, def.Path)
}

// TODO: perhaps we can move the latest tag handling into the plugin package itself? is it smelly here? if u have free time, think of the perfect abstraction for this.

// isLatest reports whether v is the literal "latest" or "latest:<prefix>".
func isLatest(v string) bool {
	return v == "latest" || strings.HasPrefix(v, "latest:")
}

// latestPrefix extracts the prefix from "latest:<prefix>" (returns "" for plain "latest").
func latestPrefix(v string) string {
	return strings.TrimPrefix(strings.TrimPrefix(v, "latest"), ":")
}

// callbackEnv builds the env map plugin scripts expect.
func callbackEnv(target Target, installPath, downloadPath string) map[string]string {
	return plugin.AsdfEnv{
		Version:      target.Version,
		InstallPath:  installPath,
		DownloadPath: downloadPath,
		Concurrency:  runtime.NumCPU(),
	}.Map()
}
