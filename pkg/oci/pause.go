package oci

import (
	"os/exec"

	"github.com/pkg/errors"
)

// Pause pauses a container
func Pause(ociBinary string, ociID string) error {
	cmd := exec.Command(ociBinary, "pause", ociID)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error pausing node %s", ociID)
	}

	return nil
}
