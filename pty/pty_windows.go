//go:build windows

package pty

import (
	"bufio"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"os/exec"
)

type Console struct {
	Stream *bufio.ReadWriter
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func (c *Console) Start() error {

	var err error
	c.cmd = exec.Command("powershell.exe", "-NoExit")
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		return err
	}

	c.stdout, err = c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	reader := transform.NewReader(c.stdout, simplifiedchinese.GBK.NewDecoder())

	c.cmd.Stderr = os.Stderr
	c.Stream = bufio.NewReadWriter(bufio.NewReader(reader), bufio.NewWriter(c.stdin))
	err = c.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Console) Stop() {
	if c.stderr != nil {
		c.stderr.Close()
	}
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.stdout != nil {
		c.stdout.Close()
	}
	c.cmd.Process.Kill()
}
