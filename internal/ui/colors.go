// Package ui handles all outputs.
package ui

import (
	"os"

	"golang.org/x/term"
)

// color is an ANSI escape code.
type color string

const (
	reset  color = "\033[0m"
	red    color = "\033[31m"
	green  color = "\033[32m"
	yellow color = "\033[33m"
	dim    color = "\033[2m"
)

// Per-stream interactive-terminal and color decisions, computed once at package init.
var (
	isStderrTTY = ttyMode(os.Stderr)
	colorStderr = isStderrTTY && !noColor()

	// compute for stdout as needed..
)

// ttyMode reports whether f is an interactive terminal where cursor escapes are appropriate.
// CI=true is treated as non-TTY so logs stay clean.
func ttyMode(f *os.File) bool {
	if os.Getenv("CI") == "true" {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

// noColor reports whether ANSI color codes should be suppressed (per no-color.org).
func noColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

// paint wraps s in the given color when on is true, otherwise returns s unchanged.
func paint(s string, c color, on bool) string {
	if !on {
		return s
	}
	return string(c) + s + string(reset)
}

// Green wraps s in green ANSI when stderr would emit color.
func Green(s string) string { return paint(s, green, colorStderr) }

// Red wraps s in red ANSI when stderr would emit color.
func Red(s string) string { return paint(s, red, colorStderr) }

// Dim wraps s in dim ANSI when stderr would emit color.
func Dim(s string) string { return paint(s, dim, colorStderr) }

// ...add as needed :)
