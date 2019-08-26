package action

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/pkg/errors"

	"github.com/medyagh/kic/pkg/runner"
)

// InstallCNI installs CNI
func InstallCNI(r runner.Cmder, subnet string) error {
	// read the manifest from the node
	var raw bytes.Buffer
	if err := r.Command("cat", "/kind/manifests/default-cni.yaml").SetStdout(&raw).Run(); err != nil {
		return errors.Wrap(err, "failed to read CNI manifest")
	}
	manifest := raw.String()

	if strings.Contains(manifest, "would you kindly template this file") {
		t, err := template.New("cni-manifest").Parse(manifest)
		if err != nil {
			return errors.Wrap(err, "failed to parse CNI manifest template")
		}
		var out bytes.Buffer
		err = t.Execute(&out, &struct {
			PodSubnet string
		}{
			PodSubnet: subnet,
		})
		if err != nil {
			return errors.Wrap(err, "failed to execute CNI manifest template")
		}
		manifest = out.String()
	}

	if err := r.Command(
		"kubectl", "create", "--kubeconfig=/etc/kubernetes/admin.conf",
		"-f", "-",
	).SetStdin(strings.NewReader(manifest)).Run(); err != nil {
		return errors.Wrap(err, "failed to apply overlay network")
	}

	return nil

}
