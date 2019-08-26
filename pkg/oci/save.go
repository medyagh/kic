package oci

import (
	"github.com/medyagh/kic/pkg/runner"
	"github.com/pkg/errors"
)

// Save saves an image archive "docker/podman save"
func Save(image, dest string) error {
	cmd := runner.Command(DefaultOCI, "save", "-o", dest, image)
	lines, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "saving image to tar failed, output %s", lines[0])
	}
	return nil
}
