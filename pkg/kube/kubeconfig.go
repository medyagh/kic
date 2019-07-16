package kube

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/medyagh/kic/pkg/exec"
	"github.com/medyagh/kic/pkg/node"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/homedir"
)

// rename generate based on /etc/...
func GenerateKubeConfig(n *node.Node, hostIP string, hostPort int32, profile string) ([]byte, error) {
	cmd := n.Command("cat", "/etc/kubernetes/admin.conf")
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubeconfig from node")
	}

	// fix the config file, swapping out the server for the forwarded localhost:port
	var buff bytes.Buffer
	for _, line := range lines {
		match := serverAddressRE.FindStringSubmatch(line)
		if len(match) > 1 {
			addr := net.JoinHostPort(hostIP, fmt.Sprintf("%d", hostPort))
			line = fmt.Sprintf("%s https://%s", match[1], addr)
		}
		buff.WriteString(line)
		buff.WriteString("\n")
	}

	return buff.Bytes(), nil
}

func WriteKubeConfig(content []byte, profile string) error {
	// copies the kubeconfig files locally in order to make the cluster
	// usable with kubectl.
	// the kubeconfig file created by kubeadm internally to the node
	// must be modified in order to use the random host port reserved
	// for the API server and exposed by the node
	configDir := filepath.Join(homedir.HomeDir(), ".kube")
	fileName := fmt.Sprintf("kic-config-%s", profile)
	kubeConfigPath := filepath.Join(configDir, fileName)

	// create the directory to contain the KUBECONFIG file.
	// 0755 is taken from client-go's config handling logic: https://github.com/kubernetes/client-go/blob/5d107d4ebc00ee0ea606ad7e39fd6ce4b0d9bf9e/tools/clientcmd/loader.go#L412
	fmt.Printf("\nexport KUBECONFIG=%s \n", kubeConfigPath)
	err := os.MkdirAll(filepath.Dir(kubeConfigPath), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create kubeconfig output directory")
	}

	return ioutil.WriteFile(kubeConfigPath, content, 0600)
}

// matches kubeconfig server entry like:
//    server: https://172.17.0.2:6443
// which we rewrite to:
//    server: https://$ADDRESS:$PORT
var serverAddressRE = regexp.MustCompile(`^(\s+server:) https://.*:\d+$`)
