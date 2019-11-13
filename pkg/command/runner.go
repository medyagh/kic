package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/medyagh/kic/pkg/assets"
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

// Command returns a human readable command string that does not induce eye fatigue
func (rr RunResult) Command() string {
	var sb strings.Builder
	sb.WriteString(rr.Args[0])
	for _, a := range rr.Args[1:] {
		if strings.Contains(a, " ") {
			sb.WriteString(fmt.Sprintf(` "%s"`, a))
			continue
		}
		sb.WriteString(fmt.Sprintf(" %s", a))
	}
	return sb.String()
}

type Runner interface {
	// RunCmd runs a cmd of exec.Cmd type. allowing user to set cmd.Stdin, cmd.Stdout,...
	// not all implementors are guaranteed to handle all the properties of cmd.
	RunCmd(cmd *exec.Cmd) (*RunResult, error)
	Copy(assets.LegacyCopyableFile) error
	Remove(assets.LegacyCopyableFile) error
}
