// Package ui handles all outputs.
package ui

import (
	"fmt"
	"os"
)

// Write writes to stdout.
func Write(out string) {
	_, _ = fmt.Fprintln(os.Stdout, out)
}

// Info writes to stderr without color.
func Info(out string) {
	fmt.Fprintln(os.Stderr, out)
}

// Error writes to stderr in red.
func Error(out string) {
	fmt.Fprintln(os.Stderr, paint(out, red, colorStderr))
}
