// Package binary handles entrypoint binary generation and resolution.
package binary

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
)

//go:embed scripts/shim.sh.in
var shimTemplateRaw string

var shimTemplate = template.Must(template.New("shim").Parse(shimTemplateRaw))

// Generate writes shim scripts for every binary the plugin exposes.
func Generate(ctx context.Context, p *plugin.Plugin, version, installPath string) error {
	paths, err := binPaths(ctx, p, version, installPath)
	if err != nil {
		return err
	}
	shimsDir := config.ShimsDir()
	if err := os.MkdirAll(shimsDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", shimsDir, err)
	}
	for _, rel := range paths {
		binDir := filepath.Join(installPath, rel)
		entries, err := os.ReadDir(binDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read %s: %w", binDir, err)
		}
		for _, e := range entries {
			info, err := e.Info()
			if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
				continue
			}
			shimPath := filepath.Join(shimsDir, e.Name())
			if err := os.WriteFile(shimPath, []byte(shimContent(p.Name, e.Name())), 0o755); err != nil {
				return fmt.Errorf("write shim %s: %w", shimPath, err)
			}
		}
	}
	return nil
}

// binPaths returns the relative bin paths the plugin exposes, defaulting to ["bin"].
func binPaths(ctx context.Context, p *plugin.Plugin, version, installPath string) ([]string, error) {
	// NOTE: Plugins are usually version-agnostic. All installations (all versions) use the same plugin as the single source of truth, but the asdf protocol allows us to pass it, and *might* be used by some plugins, idk.
	env := plugin.AsdfEnv{Version: version, InstallPath: installPath}.Map()
	stdout, _, err := p.RunCallback(ctx, "list-bin-paths", nil, env)
	if err != nil {
		if errors.Is(err, plugin.ErrCallbackNotImplemented) {
			return []string{"bin"}, nil
		}
		return nil, err
	}
	paths := strings.Fields(string(stdout))
	if len(paths) == 0 {
		return []string{"bin"}, nil
	}
	return paths, nil
}

// shimContent produces the shim script for binname owned by plugin.
func shimContent(pluginName, binname string) string {
	var buf bytes.Buffer
	_ = shimTemplate.Execute(&buf, struct {
		Plugin, BinName string
	}{pluginName, binname})
	return buf.String()
}
