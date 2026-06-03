// Package config defines project-wide configuration.
package config

import (
	"os"
	"path/filepath"
)

const dataDirFolder = ".zxcv"
const (
	installsSubdir  = "installs"
	downloadsSubdir = "downloads"
	pluginsSubdir   = "plugins"
	shimsSubdir     = "bin"
	registrySubdir  = "plugin-index"
	libexecSubdir   = "libexec"
)

// GlobalDir returns the directory where the global `.tool-versions` file lives ($HOME).
func GlobalDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// DataDir contains our root data directory ($HOME/.zxcv).
func DataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return dataDirFolder
	}
	return filepath.Join(home, dataDirFolder)
}

// InstallsDir is where installed tools live.
func InstallsDir() string { return filepath.Join(DataDir(), installsSubdir) }

// DownloadsDir is the download area, cleaned up after install.
func DownloadsDir() string { return filepath.Join(DataDir(), downloadsSubdir) }

// PluginsDir is where plugin repos and symlinks live.
func PluginsDir() string { return filepath.Join(DataDir(), pluginsSubdir) }

// ShimsDir is where all entrypoint binaries live.
func ShimsDir() string { return filepath.Join(DataDir(), shimsSubdir) }

// RegistryDir is the local clone of the asdf plugin index.
func RegistryDir() string { return filepath.Join(DataDir(), registrySubdir) }

// LibexecDir holds internal helper executables not meant to be on the user's PATH.
func LibexecDir() string { return filepath.Join(DataDir(), libexecSubdir) }
