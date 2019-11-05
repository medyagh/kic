package action

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"

	"github.com/medyagh/kic/pkg/command"
	"github.com/pkg/errors"
)

// GetDefaultCNIManifest returns the default CNI manifest
func GetDefaultCNIManifest(r command.Runner, subnet string) ([]byte, error) {
	// read the manifest from the node
	var raw bytes.Buffer
	cmd := exec.Command("cat", "/kind/manifests/default-cni.yaml")
	cmd.Stdout = &raw

	if _, err := r.RunCmd(cmd); err != nil {
		return nil, errors.Wrap(err, "failed to read CNI manifest")
	}
	manifest := raw.String()

	if !strings.Contains(manifest, "would you kindly template this file") {
		return nil, errors.New("bad default CNI template")
	}

	t, err := template.New("cni-manifest").Parse(manifest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse CNI manifest template")
	}

	var out bytes.Buffer
	err = t.Execute(&out, &struct {
		PodSubnet string
	}{
		PodSubnet: subnet,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to execute CNI manifest template")
	}

	return out.Bytes(), nil
}

// ApplyCNIManifest applies a CNI manifest
func ApplyCNIManifest(r command.Runner, manifest []byte) error {
	cmd := exec.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "-",
	)
	cmd.Stdin = bytes.NewReader(manifest)
	if _, err := r.RunCmd(cmd); err != nil {
		return errors.Wrap(err, "failed to apply overlay network")
	}
	return nil
}
