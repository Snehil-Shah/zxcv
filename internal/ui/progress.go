// Package ui handles all outputs.
package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// TODO: this works and is tailored for our use, but we'll eventually migrate to bubbletea when we decide to support more complex UIs like in-place log streaming (like the Docker CLI).

// spinnerFrames is the braille glyph cycle used while a label is pending. looks pretty, and is standard.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// progressState tracks per-label completion.
type progressState int

const (
	statePending progressState = iota
	stateDone
	stateFailed
)

// Progress renders a list of tasks with live spinners on TTYs, plain lines elsewhere.
type Progress struct {
	labels   []string
	states   []progressState
	mu       sync.Mutex
	frame    int
	ticker   *time.Ticker
	done     chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// NewProgress prints initial labels and starts the spinner if stderr is a TTY.
func NewProgress(labels []string) *Progress {
	p := &Progress{
		labels: labels,
		states: make([]progressState, len(labels)),
	}
	if isStderrTTY {
		// Ticker must exist before the initial render so formatLine picks the spinner glyph.
		p.ticker = time.NewTicker(100 * time.Millisecond)
		p.done = make(chan struct{})
		for i := range labels {
			fmt.Fprintln(os.Stderr, p.formatLine(i))
		}
		p.wg.Add(1)
		go p.spin()
	}
	return p
}

// Done marks idx as completed, auto-stops the spinner once every line has finished.
func (p *Progress) Done(idx int) { p.update(idx, stateDone) }

// Failed marks idx as failed, auto-stops the spinner once every line has finished.
func (p *Progress) Failed(idx int) { p.update(idx, stateFailed) }

// Stop halts the spinner, final states stay on screen. Idempotent and safe to call after auto-stop.
func (p *Progress) Stop() {
	p.stopOnce.Do(func() {
		if p.ticker == nil {
			return
		}
		close(p.done)
		p.ticker.Stop()
		p.wg.Wait()
	})
}

// Clear removes the rendered progress lines from the screen. No-op on non-TTY.
// Call only after the spinner has stopped (auto-stop or Stop).
func (p *Progress) Clear() {
	if p.ticker == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "\033[%dA\033[0J", len(p.labels))
}

// update sets the state for idx, re-renders the line, and auto-stops once nothing is pending.
func (p *Progress) update(idx int, s progressState) {
	p.mu.Lock()
	p.states[idx] = s
	if p.ticker != nil {
		p.renderInPlace(idx)
	} else {
		fmt.Fprintln(os.Stderr, p.formatLine(idx))
	}
	allDone := true
	for _, st := range p.states {
		if st == statePending {
			allDone = false
			break
		}
	}
	p.mu.Unlock()
	if allDone {
		p.Stop()
	}
}

// spin advances the spinner frame and re-renders pending lines until done is closed.
func (p *Progress) spin() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ticker.C:
			p.mu.Lock()
			p.frame++
			for i, s := range p.states {
				if s == statePending {
					p.renderInPlace(i)
				}
			}
			p.mu.Unlock()
		case <-p.done:
			return
		}
	}
}

// renderInPlace overwrites the line for idx using ANSI cursor save/restore.
func (p *Progress) renderInPlace(idx int) {
	upN := len(p.labels) - idx
	// Write ANSI codes for various term operations like moving the cursor, clearing shit and whatnot.
	fmt.Fprintf(os.Stderr, "\033[s\033[%dA\r\033[2K%s\033[u", upN, p.formatLine(idx))
}

// formatLine returns the rendered "<icon> <label>" string for idx based on its state.
func (p *Progress) formatLine(idx int) string {
	var icon string
	switch p.states[idx] {
	case statePending:
		if p.ticker != nil {
			icon = paint(spinnerFrames[p.frame%len(spinnerFrames)], yellow, colorStderr)
		} else {
			icon = "▸"
		}
	case stateDone:
		icon = paint("✓", green, colorStderr)
	case stateFailed:
		icon = paint("✗", red, colorStderr)
	}
	return icon + " " + p.labels[idx]
}
