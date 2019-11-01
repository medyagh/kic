package oci

import (
	"os/exec"

	"github.com/pkg/errors"
)

// Remove removes a container
func Remove(ociID string) error {
	// TODO: force remove should be an option
	cmd := exec.Command(DefaultOCI, "rm", "-f", "-v", ociID)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error removing node %s", ociID)
	}

	return nil
}
