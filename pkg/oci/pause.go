package oci

import (
	"github.com/medyagh/kic/pkg/runner"

	"github.com/pkg/errors"
)

// Pause pauses a container
func Pause(ociID string) error {
	cmd := runner.Command(DefaultOCI, "pause", ociID)
	_, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "error pausing node %s", ociID)
	}

	return nil
}
