package oci

import (
	"os/exec"
	"strings"

	"github.com/docker/machine/libmachine/state"
	"github.com/pkg/errors"
)

// Status stops a container
func Status(ociBinary string, ociID string) (state.State, error) {
	cmd := exec.Command(ociBinary, "inspect", "-f", "{{.State.Status}}", ociID)
	out, err := cmd.CombinedOutput()
	o := strings.Trim(string(out), "\n")
	s := state.Error
	if o == "running" { // TODO: parse all kind of states
		s = state.Running
	}
	if o == "exited" {
		s = state.Stopped
	}

	if o == "paused" {
		s = state.Paused
	}

	if err != nil {
		return state.Error, errors.Wrapf(err, "error stop node %s", ociID)
	}
	return s, nil
}

// SystemStatus checks if the oci container engine is running
func SystemStatus(ociBinary string, ociID string) (state.State, error) {
	_, err := exec.LookPath(ociBinary)
	if err != nil {
		return state.Error, err
	}

	err = exec.Command("docker", "info").Run()
	if err != nil {
		return state.Error, err
	}

	return state.Running, nil
}
