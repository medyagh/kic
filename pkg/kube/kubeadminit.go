package kube

import (
	"github.com/pkg/errors"

	"github.com/medyagh/kic/pkg/exec"
	"github.com/medyagh/kic/pkg/node"
)

/// RunKubeadmInit runs kubeadm init on a node
func RunKubeadmInit(node *node.Node, hostIP string, hostPort int32, profile string) ([]string, error) { // run kubeadm
	cmd := node.Command(
		// init because this is the control plane node
		"kubeadm", "init",
		"--ignore-preflight-errors=all",
		// specify our generated config file
		"--config=/kind/kubeadm.conf",
		"--skip-token-print",
		// increase verbosity for debugging
		"--v=6",
	)
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		return lines, errors.Wrap(err, "failed to init node with kubeadm")
	}

	return lines, nil
}

func RunTaint(n *node.Node) error {
	// if we are only provisioning one node, remove the master taint
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#master-isolation
	if err := n.Command(
		"kubectl", "--kubeconfig=/etc/kubernetes/admin.conf",
		"taint", "nodes", "--all", "node-role.kubernetes.io/master-",
	).Run(); err != nil {
		return errors.Wrap(err, "failed to remove master taint")
	}
	return nil
}
