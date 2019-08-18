package oci

import (
	"github.com/medyagh/kic/pkg/exec"
)

// Save saves image to dest, as in `docker save`
func Save(image, dest string) error {
	return exec.Command(DefaultOCI, "save", "-o", dest, image).Run()
}
