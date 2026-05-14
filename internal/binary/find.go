// Package binary handles entrypoint binary generation and resolution.
package binary

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
	"github.com/Snehil-Shah/zxcv/internal/resolver"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

// Resolved describes the target binary. Plugin is nil for system-version resolutions.
type Resolved struct {
	BinPath     string
	Plugin      *plugin.Plugin
	Version     string
	InstallPath string
}

// Find resolves plugin's binary at the version applicable to cwd.
func Find(ctx context.Context, pluginName, binary, cwd string) (*Resolved, error) {
	version, err := resolver.ResolveVersion(cwd, pluginName)
	if err != nil {
		return nil, err
	}
	return resolveInstalled(ctx, pluginName, binary, version)
}

// FindAt resolves plugin's binary at the given explicit version.
func FindAt(ctx context.Context, pluginName, binary, version string) (*Resolved, error) {
	return resolveInstalled(ctx, pluginName, binary, version)
}

// PluginForShim reads binary name's shim file and returns the plugin it was generated for.
func PluginForShim(name string) (*plugin.Plugin, error) {
	shimPath := filepath.Join(config.ShimsDir(), name)
	data, err := os.ReadFile(shimPath)
	if err != nil {
		return nil, fmt.Errorf("no shim for %s: %w", name, err)
	}
	pluginName, ok := parsePluginComment(string(data))
	if !ok {
		return nil, fmt.Errorf("malformed shim %s", shimPath)
	}
	return plugin.New(pluginName), nil
}

// resolveInstalled returns the path to the binary for an installed plugin@version, or the system fallback.
func resolveInstalled(ctx context.Context, pluginName, binary, version string) (*Resolved, error) {
	if version == toolversions.SystemVersion {
		path, err := findSystem(binary, config.ShimsDir())
		if err != nil {
			return nil, err
		}
		return &Resolved{BinPath: path}, nil
	}
	installPath := filepath.Join(config.InstallsDir(), pluginName, version)
	if _, err := os.Stat(installPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s %s is not installed", pluginName, version)
		}
		return nil, fmt.Errorf("stat %s: %w", installPath, err)
	}
	p := plugin.New(pluginName)
	dirs, err := executableDirs(ctx, p, version, installPath)
	if err != nil {
		return nil, fmt.Errorf("%s %s: query bin paths: %w", pluginName, version, err)
	}

	// Walk all binary dirs the plugin declares, filter to real executables, and ask the plugin to rewrite the path if it can.
	for _, dir := range dirs {
		candidate := filepath.Join(dir, binary)
		info, statErr := os.Stat(candidate)
		if statErr != nil || info.IsDir() || info.Mode()&0o111 == 0 {
			// Not the one.
			continue
		}
		final, err := execPath(ctx, p, version, binary, candidate, installPath)
		if err != nil {
			return nil, fmt.Errorf("%s %s: exec-path: %w", pluginName, version, err)
		}
		return &Resolved{
			BinPath:     final,
			Plugin:      p,
			Version:     version,
			InstallPath: installPath,
		}, nil
	}
	return nil, fmt.Errorf("%s %s: executable %s not found in %v", pluginName, version, binary, dirs)
}

// executableDirs returns the absolute dirs where the plugin's binaries may live.
func executableDirs(ctx context.Context, p *plugin.Plugin, version, installPath string) ([]string, error) {
	var dirs []string
	shimsDir := filepath.Join(p.Dir, "shims")
	if info, err := os.Stat(shimsDir); err == nil && info.IsDir() {
		// NOTE: Plugins may define their own shims directory.
		dirs = append(dirs, shimsDir)
	}
	rel, err := binPaths(ctx, p, version, installPath)
	if err != nil {
		return nil, err
	}
	for _, r := range rel {
		// Their own declared bin paths.
		dirs = append(dirs, filepath.Join(installPath, r))
	}
	return dirs, nil
}

// execPath asks the plugin's bin/exec-path (if defined) to rewrite the resolved binary path.
func execPath(ctx context.Context, p *plugin.Plugin, version, binary, candidate, installPath string) (string, error) {
	rel, err := filepath.Rel(installPath, candidate)
	if err != nil {
		return "", fmt.Errorf("rel %s: %w", candidate, err)
	}
	env := plugin.AsdfEnv{Version: version, InstallPath: installPath}.Map()
	stdout, _, err := p.RunCallback(ctx, "exec-path", []string{installPath, binary, rel}, env)
	if err != nil {
		if errors.Is(err, plugin.ErrCallbackNotImplemented) {
			return candidate, nil
		}
		return "", err
	}
	out := strings.TrimSpace(string(stdout))
	if out == "" {
		return candidate, nil
	}
	return filepath.Join(installPath, out), nil
}

// findSystem returns the binary's path on PATH excluding shimsDir.
func findSystem(name, shimsDir string) (string, error) {
	for _, dir := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		if dir == "" || dir == shimsDir {
			continue
		}
		candidate := filepath.Join(dir, name)
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
			continue
		}
		return candidate, nil
	}
	return "", fmt.Errorf("system binary %s not found", name)
}

// parsePluginComment extracts the plugin name from the shim's bash comment line.
func parsePluginComment(content string) (string, bool) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "# zxcv-plugin:") {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(line[len("# zxcv-plugin:"):]))
		if len(fields) >= 1 {
			return fields[0], true
		}
	}
	return "", false
}
