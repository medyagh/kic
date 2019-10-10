package oci

import (
	"github.com/medyagh/kic/pkg/runner"

	"github.com/pkg/errors"
)

// Remove removes a container
func Remove(ociID string) error {
	// TODO: force remove should be an option
	cmd := runner.Command(DefaultOCI, "rm", "-f", "-v", ociID)
	_, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "error removing node %s", ociID)
	}

	return nil
}
