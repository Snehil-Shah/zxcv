// Package tooldefinitions abstracts `.tool-definitions` manifest handling.
package tooldefinitions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseFile parses the `.tool-definitions` file at path.
func ParseFile(path string) (ToolDefinitions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return parse(data, filepath.Dir(path))
}

// parse converts the contents of a `.tool-definitions` file into a `ToolDefinitions` struct.
func parse(data []byte, baseDir string) (ToolDefinitions, error) {
	var out ToolDefinitions
	for lineNum, line := range strings.Split(string(data), "\n") {
		def, valid, err := parseLine(line, baseDir)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
		}
		if !valid {
			continue
		}
		if _, exists := out.Lookup(def.Name); exists {
			return nil, fmt.Errorf("line %d: duplicate tool %q", lineNum+1, def.Name)
		}
		out = append(out, def)
	}
	return out, nil
}

// parseLine parses a tool definition line from `.tool-definitions` manifest.
func parseLine(line, baseDir string) (Definition, bool, error) {
	// Strip trailing comment:
	if i := strings.IndexByte(line, '#'); i >= 0 {
		line = line[:i]
	}
	fields := strings.Fields(line)
	switch len(fields) {
	case 0:
		return Definition{}, false, nil
	case 1:
		return Definition{}, false, fmt.Errorf("tool %q has no source", fields[0])
	case 2:
		name, target := fields[0], fields[1]
		def := Definition{Name: name}
		if isURL(target) {
			def.URL = target
		} else {
			def.Path = resolvePath(target, baseDir)
		}
		return def, true, nil
	default:
		return Definition{}, false, fmt.Errorf("tool %q: expected 2 fields, got %d", fields[0], len(fields))
	}
}

// isURL is a quick and dirty check to see if we have a URL.
func isURL(s string) bool {
	if strings.Contains(s, "://") {
		return true
	}
	if i := strings.Index(s, "@"); i > 0 && strings.Contains(s[i+1:], ":") {
		return true
	}
	return false
}

// resolvePath returns an absolute path for p, resolving relative to `baseDir` if necessary.
func resolvePath(p, baseDir string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(baseDir, p)
}
