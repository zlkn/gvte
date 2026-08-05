package pty

import (
	"io"
)

type Winsize struct {
	Rows int
	Cols int
	X    int
	Y    int
}

type PTY interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
	Fd() uintptr
}
