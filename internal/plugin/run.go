// Package plugin manages asdf plugins.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Snehil-Shah/zxcv/internal/runner"
)

// ErrCallbackNotImplemented signals the requested callback file is absent.
var ErrCallbackNotImplemented = errors.New("callback not implemented")

// RunCallback runs the plugin's `bin/<callback>` as a subprocess.
func (p *Plugin) RunCallback(ctx context.Context, callback string, args []string, env map[string]string) (stdout, stderr []byte, err error) {
	script := filepath.Join(p.Dir, "bin", callback)
	if _, statErr := os.Stat(script); statErr != nil {
		if errors.Is(statErr, fs.ErrNotExist) {
			return nil, nil, ErrCallbackNotImplemented
		}
		return nil, nil, fmt.Errorf("stat %s: %w", script, statErr)
	}
	cmd := fmt.Sprintf("%q", script)
	for _, a := range args {
		cmd += fmt.Sprintf(" %q", a)
	}
	return runner.Run(ctx, cmd, env)
}

// ExecWithEnv exec's binPath, sourcing the plugin's `bin/exec-env` first when present.
func (p *Plugin) ExecWithEnv(binPath string, args []string, extras map[string]string) error {
	script := filepath.Join(p.Dir, "bin", "exec-env")
	if _, statErr := os.Stat(script); statErr != nil {
		if errors.Is(statErr, fs.ErrNotExist) {
			return runner.Exec(binPath, args, extras)
		}
		return fmt.Errorf("stat %s: %w", script, statErr)
	}
	return runner.ExecBash(binPath, fmt.Sprintf(`source %q && exec %q "$@"`, script, binPath), args, extras)
}
