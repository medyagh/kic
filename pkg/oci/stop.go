package oci

import (
	"os/exec"

	"github.com/pkg/errors"
)

// Stop stops a container
func Stop(ociID string) error {
	cmd := exec.Command(DefaultOCI, "stop", ociID)
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "error stop node %s", ociID)
	}

	return nil
}
