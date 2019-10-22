package action

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/pkg/errors"

	"github.com/medyagh/kic/pkg/runner"
)

// GetDefaultCNIManifest returns the default CNI manifest
func GetDefaultCNIManifest(r runner.Cmder, subnet string) ([]byte, error) {
	// read the manifest from the node
	var raw bytes.Buffer
	if err := r.Command("cat", "/kind/manifests/default-cni.yaml").SetStdout(&raw).Run(); err != nil {
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
func ApplyCNIManifest(r runner.Cmder, manifest []byte) error {
	if err := r.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "-",
	).SetStdin(bytes.NewReader(manifest)).Run(); err != nil {
		return errors.Wrap(err, "failed to apply overlay network")
	}

	return nil
}
