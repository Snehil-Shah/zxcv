// Package toolversions abstracts `.tool-versions` manifest handling.
package toolversions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// WriteEntry sets/appends name's version in the `.tool-versions` file at path.
func WriteEntry(path, name, version string) error {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read %s: %w", path, err)
	}
	out := updateOrAppend(string(data), name, version)
	return os.WriteFile(path, []byte(out), 0o644)
}

// updateOrAppend rewrites or adds the entry for name in content.
func updateOrAppend(content, name, version string) string {
	if content == "" {
		return name + " " + version + "\n"
	}

	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	found := false
	for i, line := range lines {
		if !found && lineMatchesTool(line, name) {
			lines[i] = rewriteLine(line, name, version)
			found = true
		}
	}
	if !found {
		lines = append(lines, name+" "+version)
	}
	return strings.Join(lines, "\n") + "\n"
}

// lineMatchesTool reports whether line declares an entry for name.
func lineMatchesTool(line, name string) bool {
	code := line
	if i := strings.IndexByte(code, '#'); i >= 0 {
		code = code[:i]
	}
	fields := strings.Fields(code)
	return len(fields) > 0 && fields[0] == name
}

// rewriteLine replaces the tool/version portion while preserving any trailing comment.
func rewriteLine(line, name, version string) string {
	if i := strings.IndexByte(line, '#'); i >= 0 {
		return name + " " + version + " " + line[i:]
	}
	return name + " " + version
}
