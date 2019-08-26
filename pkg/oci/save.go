package oci

import (
	"sigs.k8s.io/kind/pkg/exec"
)

// Save saves an image archive "docker/podman save"
func Save(image, dest string) error {
	cmd := exec.Command(DefaultOCI, "save", "-o", dest, image)
	_, err := exec.CombinedOutputLines(cmd)
	return err
}
