package pty

import (
	"io"
)

// Winsize represents terminal window dimensions.
type Winsize struct {
	Rows int
	Cols int
	X    int
	Y    int
}

// PTY represents an operating system pseudo-terminal session.
type PTY interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
	Fd() uintptr
}
