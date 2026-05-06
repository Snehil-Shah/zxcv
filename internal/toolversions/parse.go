// Package toolversions abstracts `.tool-versions` manifest handling.
package toolversions

import (
	"fmt"
	"os"
	"strings"
)

// ParseFile parses the `.tool-versions` file at path.
func ParseFile(path string) (ToolVersions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return parse(data)
}

// parse converts the contents of a `.tool-versions` file into a `ToolVersions` struct.
func parse(data []byte) (ToolVersions, error) {
	var out ToolVersions
	for lineNum, line := range strings.Split(string(data), "\n") {
		tool, valid, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
		}
		if !valid {
			continue
		}
		if _, exists := out.Lookup(tool.Name); exists {
			return nil, fmt.Errorf("line %d: duplicate tool %q", lineNum+1, tool.Name)
		}
		out = append(out, tool)
	}
	return out, nil
}

// parseLine parses a single tool entry. Exactly one version per tool is allowed.
func parseLine(line string) (Tool, bool, error) {
	// Strip trailing comment:
	if i := strings.IndexByte(line, '#'); i >= 0 {
		line = line[:i]
	}
	fields := strings.Fields(line)
	switch len(fields) {
	case 0:
		return Tool{}, false, nil
	case 1:
		return Tool{}, false, fmt.Errorf("tool %q has no version", fields[0])
	case 2:
		return Tool{Name: fields[0], Version: fields[1]}, true, nil
	default:
		return Tool{}, false, fmt.Errorf("tool %q: only one version per line allowed (got %d)", fields[0], len(fields)-1)
	}
}
