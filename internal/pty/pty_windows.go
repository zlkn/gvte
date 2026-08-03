//go:build windows

package pty

import (
	"fmt"
	"os"
	"os/exec"
)

// WindowsPTY implements PTY interface using Windows ConPTY API.
type WindowsPTY struct {
	hPC uintptr
	in  *os.File
	out *os.File
}

// Start launches a command attached to a new ConPTY session on Windows.
func Start(cmd *exec.Cmd, sz *Winsize) (PTY, error) {
	// Stub for CreatePseudoConsole / ConPTY initialization on Windows
	return nil, fmt.Errorf("Windows ConPTY implementation pending syscall setup")
}

func (p *WindowsPTY) Read(b []byte) (int, error) {
	if p.out == nil {
		return 0, os.ErrInvalid
	}
	return p.out.Read(b)
}

func (p *WindowsPTY) Write(b []byte) (int, error) {
	if p.in == nil {
		return 0, os.ErrInvalid
	}
	return p.in.Write(b)
}

func (p *WindowsPTY) Close() error {
	if p.in != nil {
		_ = p.in.Close()
	}
	if p.out != nil {
		_ = p.out.Close()
	}
	return nil
}

func (p *WindowsPTY) Resize(cols, rows int) error {
	// Stub for ResizePseudoConsole call
	return nil
}

func (p *WindowsPTY) Fd() uintptr {
	return 0
}
