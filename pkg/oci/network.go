package oci

import (
	"strings"

	"github.com/medyagh/kic/pkg/exec"
)

// NetworkInspect displays detailed information on one or more networks
func NetworkInspect(networkNames []string, format string) ([]string, error) {
	cmd := exec.Command("docker", "network", "inspect",
		"-f", format,
		strings.Join(networkNames, " "),
	)
	return exec.CombinedOutputLines(cmd)
}
