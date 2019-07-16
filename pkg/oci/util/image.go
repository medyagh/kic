package util

import (
	"github.com/medyagh/kic/pkg/exec"
	"github.com/pkg/errors"
)

// ImageInspect return low-level information on containers images
func ImageInspect(containerNameOrID, format string) ([]string, error) {
	cmd := exec.Command("docker", "image", "inspect",
		"-f", format,
		containerNameOrID, // ... against the container
	)

	return exec.CombinedOutputLines(cmd)
}

// ImageID return the Id of the container image
func ImageID(containerNameOrID string) (string, error) {
	lines, err := ImageInspect(containerNameOrID, "{{ .Id }}")
	if err != nil {
		return "", err
	}
	if len(lines) != 1 {
		return "", errors.Errorf("Docker image ID should only be one line, got %d lines", len(lines))
	}
	return lines[0], nil
}
