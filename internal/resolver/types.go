// Package resolver resolves the applicable manifest files.
package resolver

import (
	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

// Resolution is a single resolved tool, the file that supplied it, and any accompanying definition.
type Resolution struct {
	toolversions.Tool
	Source     string                      // path to the `.tool-versions` file
	Definition *tooldefinitions.Definition // sibling `.tool-definitions` entry, if present
}
