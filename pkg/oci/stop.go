package oci

import (
	"github.com/medyagh/kic/pkg/runner"

	"github.com/pkg/errors"
)

// Stop stops a container
func Stop(ociID string) error {
	cmd := runner.Command(DefaultOCI, "stop", ociID)
	_, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "error stop node %s", ociID)
	}

	return nil
}
