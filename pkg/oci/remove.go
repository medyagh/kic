package oci

import (
	"os/exec"

	"github.com/pkg/errors"
)

// Remove removes a container
func Remove(ociBinary string, ociID string) error {
	// TODO: force remove should be an option
	cmd := exec.Command(ociBinary, "rm", "-f", "-v", ociID)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error removing node %s", ociID)
	}

	return nil
}
