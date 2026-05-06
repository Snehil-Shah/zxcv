// Package plugin manages asdf plugins.
package plugin

import "strconv"

// AsdfEnv builds the ASDF_* env vars asdf plugin scripts expect. Zero-valued fields are omitted.
type AsdfEnv struct {
	Version      string
	InstallPath  string
	PluginPath   string // ASDF_PLUGIN_PATH (exec-env)
	DownloadPath string // ASDF_DOWNLOAD_PATH (download/install callbacks)
	Concurrency  int    // ASDF_CONCURRENCY (install callback)
}

// Map returns the env map; empty/zero fields are skipped.
func (a AsdfEnv) Map() map[string]string {
	m := map[string]string{}
	if a.Version != "" {
		m["ASDF_INSTALL_TYPE"] = "version"
		m["ASDF_INSTALL_VERSION"] = a.Version
	}
	if a.InstallPath != "" {
		m["ASDF_INSTALL_PATH"] = a.InstallPath
	}
	if a.PluginPath != "" {
		m["ASDF_PLUGIN_PATH"] = a.PluginPath
	}
	if a.DownloadPath != "" {
		m["ASDF_DOWNLOAD_PATH"] = a.DownloadPath
	}
	if a.Concurrency > 0 {
		m["ASDF_CONCURRENCY"] = strconv.Itoa(a.Concurrency)
	}
	return m
}
