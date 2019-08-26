package oci

import (
	"os"

	"github.com/medyagh/kic/pkg/runner"
	"github.com/pkg/errors"
)

// Copy copies a file/folder into container
func Copy(source, dest string) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return errors.Wrapf(err, "error source %s does not exist", source)
	}

	cmd := runner.Command(DefaultOCI, "cp", source, dest)
	_, err := runner.CombinedOutputLines(cmd)
	if err != nil {
		return errors.Wrapf(err, "error copying %s into node", source)
	}
	return nil
}
