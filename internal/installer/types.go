// Package installer implements the install and uninstall pipelines.
package installer

import (
	"sync"

	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

// Target describes one tool to install.
type Target struct {
	toolversions.Tool
	Definition *tooldefinitions.Definition // if nil, resolve via registry
}

// Result is the outcome of installing a Target.
type Result struct {
	Target Target
	Stdout []byte
	Stderr []byte
	Err    error
}

// Installer runs installation of targets.
type Installer struct {
	OnComplete func(idx int, r Result) // fired serially as each target finishes
	progressMu sync.Mutex              // serializes OnComplete calls
}

// New returns a fresh Installer.
func New() *Installer {
	return &Installer{}
}
