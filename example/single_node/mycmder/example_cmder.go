package mycmder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/medyagh/kic/pkg/command"
	"github.com/pkg/errors"
	"k8s.io/klog"

	"golang.org/x/crypto/ssh/terminal"
)

// New creates a new implementor of runner
func New(containerNameOrID string, oci string) command.Runner {
	return &containerCmder{
		nameOrID: containerNameOrID,
		ociBin:   oci,
	}
}

type containerCmder struct {
	nameOrID string
	ociBin   string
}

func (c *containerCmder) RunCmd(cmd *exec.Cmd) (*command.RunResult, error) {
	args := []string{
		"exec",
		// run with privileges so we can remount etc..
		"--privileged",
	}
	if cmd.Stdin != nil {
		args = append(args,
			"-i", // interactive so we can supply input
		)
	}
	// if the command is hooked up to the processes's output we want a tty
	if isTerminal(cmd.Stderr) || isTerminal(cmd.Stdout) {
		args = append(args,
			"-t",
		)
	}
	// set env
	for _, env := range cmd.Env {
		args = append(args, "-e", env)
	}
	// specify the container and command, after this everything will be
	// args the the command in the container rather than to docker
	args = append(
		args,
		c.nameOrID, // ... against the container
	)

	args = append(
		args,
		cmd.Args...,
	)
	cmd2 := exec.Command("docker", args...)
	cmd2.Stdin = cmd.Stdin
	cmd2.Stdout = cmd.Stdout
	cmd2.Stderr = cmd.Stderr
	cmd2.Env = cmd.Env

	rr := &command.RunResult{Args: cmd.Args}

	var outb, errb io.Writer
	if cmd2.Stdout == nil {
		var so bytes.Buffer
		outb = io.MultiWriter(&so, &rr.Stdout)
	} else {
		outb = io.MultiWriter(cmd2.Stdout, &rr.Stdout)
	}

	if cmd2.Stderr == nil {
		var se bytes.Buffer
		errb = io.MultiWriter(&se, &rr.Stderr)
	} else {
		errb = io.MultiWriter(cmd2.Stderr, &rr.Stderr)
	}

	cmd2.Stdout = outb
	cmd2.Stderr = errb

	start := time.Now()

	err := cmd2.Run()
	elapsed := time.Since(start)
	if err == nil {
		// Reduce log spam
		if elapsed > (1 * time.Second) {
			klog.Infof("(ExecRunner) Done: %v: (%s)", cmd2.Args, elapsed)
		}
	} else {
		if exitError, ok := err.(*exec.ExitError); ok {
			rr.ExitCode = exitError.ExitCode()
		}
		fmt.Printf("(ExecRunner) Non-zero exit: %v: %v (%s)\n", cmd2.Args, err, elapsed)
		fmt.Printf("(ExecRunner) Output:\n %q \n", rr.Output())
		err = errors.Wrapf(err, "command failed: %s", cmd2.Args)
	}
	return rr, err

}

// IsTerminal returns true if the writer w is a terminal
func isTerminal(w io.Writer) bool {
	if v, ok := (w).(*os.File); ok {
		return terminal.IsTerminal(int(v.Fd()))
	}
	return false
}
