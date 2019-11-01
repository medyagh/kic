package runner

import (
	"io"
	osexec "os/exec"
)

// LocalCmd wraps os/runner.Cmd, implementing the runner.Cmd interface
type LocalCmd struct {
	*osexec.Cmd
}

var _ Cmd = &LocalCmd{}

// LocalCmder is a factory for LocalCmd, implementing Cmder
type LocalCmder struct{}

var _ Cmder = &LocalCmder{}

// Command returns a new runner.Cmd backed by Cmd
func (c *LocalCmder) Command(name string, arg ...string) Cmd {
	return &LocalCmd{
		Cmd: osexec.Command(name, arg...),
	}
}

// SetStdin sets stdin
func (cmd *LocalCmd) SetStdin(r io.Reader) Cmd {
	cmd.Stdin = r
	return cmd
}

// SetStdout set stdout
func (cmd *LocalCmd) SetStdout(w io.Writer) Cmd {
	cmd.Stdout = w
	return cmd
}

// SetStderr sets stderr
func (cmd *LocalCmd) SetStderr(w io.Writer) Cmd {
	cmd.Stderr = w
	return cmd
}

// Run runs
func (cmd *LocalCmd) Run() error {
	return cmd.Cmd.Run()
}
