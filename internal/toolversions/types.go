// Package toolversions abstracts `.tool-versions` manifest handling.
package toolversions

// SystemVersion is the special sentinel that exec's the binary on PATH outside our shims.
const SystemVersion = "system"

// Tool represents a `.tool-versions` entry.
type Tool struct {
	Name    string
	Version string
}

// ToolVersions represents a parsed `.tool-versions` file.
type ToolVersions []Tool

// Lookup returns the tool entry for name and whether it was found.
func (f ToolVersions) Lookup(name string) (Tool, bool) {
	for _, e := range f {
		if e.Name == name {
			return e, true
		}
	}
	return Tool{}, false
}
