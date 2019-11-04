package action

import (
	"os/exec"

	"github.com/pkg/errors"

	"github.com/medyagh/kic/pkg/command"
)

// RunKubeadmInit runs kubeadm init on a node
func RunKubeadmInit(r command.Runner, kubeadmCfgPath, profile string) error { // run kubeadm
	cmd := exec.Command(
		// init because this is the control plane node
		"kubeadm", "init",
		"--ignore-preflight-errors=all",
		// specify our generated config file
		"--config="+kubeadmCfgPath,
		"--skip-token-print",
		// increase verbosity for debugging
		"--v=6",
	)
	_, err := r.RunCmd(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to init node with kubeadm")
	}

	return nil
}

// RemoveMasterTaint removes the master node taint.
// This allows pods to be scheduled on the master node.
func RemoveMasterTaint(r command.Runner) error {
	// if we are only provisioning one node, remove the master taint
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#master-isolation
	cmd := exec.Command(
		"kubectl", "--kubeconfig=/etc/kubernetes/admin.conf",
		"taint", "nodes", "--all", "node-role.kubernetes.io/master-",
	)

	if _, err := r.RunCmd(cmd); err != nil {
		return errors.Wrap(err, "failed to remove master taint")
	}
	return nil
}
