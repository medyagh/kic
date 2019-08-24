package node

import (
	"fmt"

	"github.com/medyagh/kic/pkg/config/cri"
	"github.com/medyagh/kic/pkg/oci"
	"github.com/medyagh/kic/pkg/runner"
)

const (
	// Docker default bridge network is named "bridge" (https://docs.docker.com/network/bridge/#use-the-default-bridge-network)
	defaultNetwork  = "bridge"
	httpProxy       = "HTTP_PROXY"
	httpsProxy      = "HTTPS_PROXY"
	noProxy         = "NO_PROXY"
	ClusterLabelKey = "io.k8s.sigs.kic.cluster" // ClusterLabelKey is applied to each node docker container for identification
	NodeRoleKey     = "io.k8s.sigs.kic.role"
)

func CreateNode(name, image, clusterLabel, role string, mounts []cri.Mount, portMappings []cri.PortMapping, cmder runner.Cmder, extraArgs ...string) (*Node, error) {
	runArgs := []string{
		"-d", // run the container detached
		"-t", // allocate a tty for entrypoint logs
		// running containers in a container requires privileged
		// NOTE: we could try to replicate this with --cap-add, and use less
		// privileges, but this flag also changes some mounts that are necessary
		// including some ones docker would otherwise do by default.
		// for now this is what we want. in the future we may revisit this.
		"--privileged",
		"--security-opt", "seccomp=unconfined", // also ignore seccomp
		"--tmpfs", "/tmp", // various things depend on working /tmp
		"--tmpfs", "/run", // systemd wants a writable /run
		// some k8s things want /lib/modules
		"-v", "/lib/modules:/lib/modules:ro",
		"--hostname", name, // make hostname match container name
		"--name", name, // ... and set the container name
		// label the node with the cluster ID
		"--label", clusterLabel,
		// label the node with the role ID
		"--label", fmt.Sprintf("%s=%s", NodeRoleKey, role),
	}

	// pass proxy environment variables to be used by node's docker deamon
	proxyDetails, err := getProxyDetails()
	if err != nil || proxyDetails == nil {
		return nil, fmt.Errorf("proxy setup error : %v", err)
	}
	for key, val := range proxyDetails.Envs {
		runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", key, val))
	}

	// adds node specific args
	runArgs = append(runArgs, extraArgs...)

	if oci.UsernsRemap() {
		// We need this argument in order to make this command work
		// in systems that have userns-remap enabled on the docker daemon
		runArgs = append(runArgs, "--userns=host")
	}

	_, err = oci.CreateContainer(
		image,
		oci.WithRunArgs(runArgs...),
		oci.WithMounts(mounts),
		oci.WithPortMappings(portMappings),
	)

	// we should return a handle so the caller can clean it up
	node := FromName(name)
	node.cmder = cmder
	if err != nil {
		return node, fmt.Errorf("docker run error %v", err)
	}

	return node, nil
}

// CreateControlPlaneNode creates a contol-plane node
// and gets ready for exposing the the API server
func CreateControlPlaneNode(name, image, clusterLabel, listenAddress string, port int32, mounts []cri.Mount, portMappings []cri.PortMapping, cmder runner.Cmder) (node *Node, err error) {
	// add api server port mapping
	portMappingsWithAPIServer := append(portMappings, cri.PortMapping{
		ListenAddress: listenAddress,
		HostPort:      port,
		ContainerPort: 6443,
	})
	node, err = CreateNode(
		name, image, clusterLabel, "control-plane", mounts, portMappingsWithAPIServer, cmder,
		// publish selected port for the API server
		"--expose", fmt.Sprintf("%d", port),
	)
	if err != nil {
		return node, err
	}

	// stores the port mapping into the node internal state
	node.cache.set(func(cache *nodeCache) {
		cache.ports = map[int32]int32{6443: port}
	})
	return node, nil
}

// FromName creates a node handle from the node' Name
func FromName(name string) *Node {
	return &Node{
		name:  name,
		cache: &nodeCache{},
	}
}
