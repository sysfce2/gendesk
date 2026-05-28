package main

import "github.com/xyproto/vt"

// newSilentOutput creates a TextOutput with color and output disabled, for use in tests
func newSilentOutput() *vt.TextOutput {
	return vt.NewTextOutput(false, false)
}
