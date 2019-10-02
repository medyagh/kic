package oci

import (
	"github.com/medyagh/kic/pkg/runner"

	"github.com/pkg/errors"
)

// Delete removes a container
func Delete(ociID string) error {
	// TODO: force remove should be an option
	cmd := runner.Command(DefaultOCI, "rm", "-f", "-v", ociID)
	_, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "error deleting node %s", ociID)
	}

	return nil
}
