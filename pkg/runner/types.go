package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunResult holds the results of a Runner
type RunResult struct {
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
	Args     []string // the args that was passed to Runner
}

// Output returns human-readable output for an execution result
func (rr RunResult) Output() string {
	var sb strings.Builder
	if rr.Stdout.Len() > 0 {
		sb.WriteString(fmt.Sprintf("-- stdout --\n%s\n-- /stdout --", rr.Stdout.Bytes()))
	}
	if rr.Stderr.Len() > 0 {
		sb.WriteString(fmt.Sprintf("\n** stderr ** \n%s\n** /stderr **", rr.Stderr.Bytes()))
	}
	return sb.String()
}

type Runner interface {
	// RunCmd runs a cmd of exec.Cmd type. allowing user to set cmd.Stdin, cmd.Stdout,...
	// not all implementors are guaranteed to handle all the properties of cmd.
	RunCmd(cmd *exec.Cmd) (*RunResult, error)
}
