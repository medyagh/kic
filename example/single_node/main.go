package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/medyagh/kic/pkg/kube"
	"github.com/medyagh/kic/pkg/node"
	"github.com/medyagh/kic/pkg/node/cri"
	"github.com/medyagh/kic/pkg/oci"
	"github.com/phayes/freeport"

	"k8s.io/klog"
)

func main() {
	profile := flag.String("profile", "p1", "profile name")
	delete := flag.Bool("delete", false, "to delete")
	start := flag.Bool("start", false, "to start")
	hostIP := flag.String("host-ip", "127.0.0.1", "node's ip")

	flag.Parse()
	p, err := freeport.GetFreePort()
	hostPort := int32(p)
	if err != nil {
		log.Fatal(err)
	}

	ns := newNodeSpec(*profile, *hostIP, hostPort)

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		ns.Delete()

	}

	if *start {
		fmt.Printf("Starting on port %d\n ", hostPort)

		img := "kindest/node:v1.15.0"
		err := oci.PullIfNotPresent(img)
		if err != nil {
			klog.Errorf("Error pulling image %s", img)
		}

		// create node
		node, _ := ns.Create("kic.cluster" + *profile)

		ip, _, _ := node.IP()
		kubeADMCFG, _ := kube.GetMagicConfig(ip, *profile, "v1.15.0")

		// copy the config to the node
		if err := node.WriteFile("/kind/kubeadm.conf", kubeADMCFG); err != nil {
			klog.Errorf("failed to copy kubeadm config to node : %v", err)
		}

		kube.RunKubeadmInit(node, *hostIP, hostPort, *profile)
		kube.RunTaint(node)
		c, _ := kube.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container
		kube.WriteKubeConfig(c, *profile)
		kube.InstallCNI(node, "10.244.0.0/16")

	}

}

func newNodeSpec(profile string, hostIP string, hostPort int32) *node.Spec {
	return &node.Spec{
		Name:              profile + "control-plane",
		Image:             "kindest/node:v1.15.0@sha256:b4d092fd2b507843dd096fe6c85d06a27a0cbd740a0b32a880fe61aba24bb478",
		Role:              "control-plane",
		ExtraMounts:       []cri.Mount{},
		ExtraPortMappings: []cri.PortMapping{},
		APIServerAddress:  hostIP,
		APIServerPort:     hostPort,
		IPv6:              false,
	}
}
