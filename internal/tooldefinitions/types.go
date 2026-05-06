// Package tooldefinitions abstracts `.tool-definitions` manifest handling.
package tooldefinitions

// Definition represents a `.tool-definitions` entry.
type Definition struct {
	Name string
	URL  string // set for remote git plugins
	Path string // set for local plugins
}

// ToolDefinitions represents a parsed `.tool-definitions` file.
type ToolDefinitions []Definition

// Lookup returns the definition for name, or false if the tool is absent.
func (d ToolDefinitions) Lookup(name string) (Definition, bool) {
	for _, e := range d {
		if e.Name == name {
			return e, true
		}
	}
	return Definition{}, false
}
