//go:build !windows

package pty

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Console struct {
	Stream *os.File
	cmd    *exec.Cmd
}

func (c *Console) Start() error {
	var err error
	c.cmd = exec.Command("bash")
	c.Stream, err = pty.Start(c.cmd)
	return err
}

func (c *Console) Stop() {
	c.Stream.Close()
	c.cmd.Process.Kill()
}
