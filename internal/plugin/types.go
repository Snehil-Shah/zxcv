// Package plugin manages asdf plugins.
package plugin

import (
	"path/filepath"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// Plugin represents an asdf plugin.
type Plugin struct {
	Name string
	Dir  string
}

// New returns a Plugin for name rooted under config.PluginsDir().
func New(name string) *Plugin {
	return &Plugin{
		Name: name,
		Dir:  filepath.Join(config.PluginsDir(), name),
	}
}
