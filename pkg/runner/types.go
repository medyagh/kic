package runner

import (
	"bytes"
	"io"
	"os/exec"
)

// RunResult holds the results of a Runner
type RunResult struct {
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
	Args     []string // the args that was passed to Runner
}

type Runner interface {
	// RunCmd runs a cmd of exec.Cmd type. allowing user to set cmd.Stdin, cmd.Stdout,...
	// not all implementors are guaranteed to handle all the properties of cmd.
	RunCmd(cmd *exec.Cmd) (*RunResult, error)
}

// Cmd abstracts over running a command somewhere, this is useful for testing
type Cmd interface {
	Run() error
	// Each entry should be of the form "key=value"
	SetStdin(io.Reader) Cmd
	SetStdout(io.Writer) Cmd
	SetStderr(io.Writer) Cmd
}

// Cmder abstracts over creating commands
type Cmder interface {
	// command, args..., just like os/runner.Cmd
	Command(string, ...string) Cmd
}
