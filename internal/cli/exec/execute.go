// Package exec implements the `zxcv exec` command.
package exec

import (
	"github.com/Snehil-Shah/zxcv/internal/binary"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
	"github.com/Snehil-Shah/zxcv/internal/runner"
)

// execBinary replaces the current process with the resolved binary.
func execBinary(r *binary.Resolved, args []string) error {
	if r.Plugin == nil {
		// System. Just exec the binary.
		return runner.Exec(r.BinPath, args, nil)
	}
	extras := plugin.AsdfEnv{
		Version:     r.Version,
		InstallPath: r.InstallPath,
		PluginPath:  r.Plugin.Dir,
	}.Map()
	return r.Plugin.ExecWithEnv(r.BinPath, args, extras)
}
