package oci

import (
	"github.com/medyagh/kic/pkg/exec"
)

// Inspect return low-level information on containers
func Inspect(containerNameOrID, format string) ([]string, error) {
	cmd := exec.Command(DefaultOCI, "inspect",
		"-f", format,
		containerNameOrID, // ... against the "node" container
	)

	return exec.CombinedOutputLines(cmd)
}
