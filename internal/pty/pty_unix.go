//go:build !windows

package pty

import (
	"fmt"
	"os"
	"os/exec"
)

type UnixPTY struct {
	master *os.File
	slave  *os.File
	cmd    *exec.Cmd
}

func Start(cmd *exec.Cmd, sz *Winsize) (PTY, error) {
	//TODO Stub for Unix termios / master-slave PTY initialization
	return nil, fmt.Errorf("Unix PTY implementation pending termios setup")
}

func (p *UnixPTY) Read(b []byte) (int, error) {
	if p.master == nil {
		return 0, os.ErrInvalid
	}
	return p.master.Read(b)
}

func (p *UnixPTY) Write(b []byte) (int, error) {
	if p.master == nil {
		return 0, os.ErrInvalid
	}
	return p.master.Write(b)
}

func (p *UnixPTY) Close() error {
	if p.master != nil {
		_ = p.master.Close()
	}
	if p.slave != nil {
		_ = p.slave.Close()
	}
	return nil
}

func (p *UnixPTY) Resize(cols, rows int) error {
	//TODO Stub for TIOCSWINSZ ioctl
	return nil
}

func (p *UnixPTY) Fd() uintptr {
	if p.master == nil {
		return 0
	}
	return p.master.Fd()
}
