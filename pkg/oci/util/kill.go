package util

import (
	"github.com/medyagh/kic/pkg/exec"
)

// Kill sends the named signal to the container
func Kill(signal, containerNameOrID string) error {
	cmd := exec.Command(
		"docker", "kill",
		"-s", signal,
		containerNameOrID,
	)
	return cmd.Run()
}
