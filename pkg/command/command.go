package command

import (
	"bytes"
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
