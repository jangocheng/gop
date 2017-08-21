// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	// GlobalCommandArgs global command args for external package setting
	GlobalCommandArgs []string
)

// ErrExecTimeout error of executing timed out
type ErrExecTimeout struct {
	Timeout time.Duration // seconds
	Command *Command
}

// Error implement error interface
func (e ErrExecTimeout) Error() string {
	return "timeout failed"
}

// ConcatenateError describes the error of concate
type ConcatenateError struct {
	Err    error
	Reason string
}

// Error implement error interface
func (e ConcatenateError) Error() string {
	return "concatenate error"
}

// Command represents a command with its subcommands or arguments.
type Command struct {
	name string
	args []string
	Env  []string
}

// String implement stringer interface
func (c *Command) String() string {
	if len(c.args) == 0 {
		return c.name
	}
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

// NewCommand creates and returns a new Git Command based on given command and arguments.
func NewCommand(args ...string) *Command {
	return &Command{
		name: "go",
		args: append(GlobalCommandArgs, args...),
	}
}

// AddArguments adds new argument(s) to the command.
func (c *Command) AddArguments(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

// RunInDirTimeoutPipeline executes the command in given directory with given timeout,
// it pipes stdout and stderr to given io.Writer.
func (c *Command) RunInDirTimeoutPipeline(timeout time.Duration, dir string, stdout, stderr io.Writer) error {
	if timeout == -1 {
		timeout = 3 * time.Minute
	}

	cmd := exec.Command(c.name, c.args...)
	cmd.Dir = dir
	cmd.Env = c.Env
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); err != nil {
				return fmt.Errorf("fail to kill process: %v", err)
			}
		}

		<-done
		return ErrExecTimeout{timeout, c}
	case err = <-done:
	}

	return err
}

// RunInDirTimeout executes the command in given directory with given timeout,
// and returns stdout in []byte and error (combined with stderr).
func (c *Command) RunInDirTimeout(timeout time.Duration, dir string) ([]byte, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if err := c.RunInDirTimeoutPipeline(timeout, dir, stdout, stderr); err != nil {
		return nil, ConcatenateError{err, stderr.String()}
	}

	if stdout.Len() > 0 {
		log.Printf("stdout:\n%s\n", stdout.Bytes()[:1024])
	}
	return stdout.Bytes(), nil
}

// RunInDirPipeline executes the command in given directory,
// it pipes stdout and stderr to given io.Writer.
func (c *Command) RunInDirPipeline(dir string, stdout, stderr io.Writer) error {
	return c.RunInDirTimeoutPipeline(-1, dir, stdout, stderr)
}

// RunInDirBytes executes the command in given directory
// and returns stdout in []byte and error (combined with stderr).
func (c *Command) RunInDirBytes(dir string) ([]byte, error) {
	return c.RunInDirTimeout(-1, dir)
}

// RunInDir executes the command in given directory
// and returns stdout in string and error (combined with stderr).
func (c *Command) RunInDir(dir string) (string, error) {
	stdout, err := c.RunInDirTimeout(-1, dir)
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// RunTimeout executes the command in default working directory with given timeout,
// and returns stdout in string and error (combined with stderr).
func (c *Command) RunTimeout(timeout time.Duration) (string, error) {
	stdout, err := c.RunInDirTimeout(timeout, "")
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// Run executes the command in default working directory
// and returns stdout in string and error (combined with stderr).
func (c *Command) Run() (string, error) {
	return c.RunTimeout(-1)
}
