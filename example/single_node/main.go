package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/medyagh/kic/example/single_node/mycmder"
	"github.com/medyagh/kic/pkg/image"
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
	kubeVersion := flag.String("kubernetes-version", "v1.15.0", "kuberentes version")

	flag.Parse()
	p, err := freeport.GetFreePort()
	hostPort := int32(p)
	if err != nil {
		log.Fatal(err)
	}

	imgSha, _ := image.NameForVersion(*kubeVersion)

	ns := newNodeSpec(*profile, imgSha, *hostIP, hostPort)

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		ns.Delete()

	}

	if *start {
		fmt.Printf("Starting on port %d\n ", hostPort)
		err := oci.PullIfNotPresent(imgSha)
		if err != nil {
			klog.Errorf("Error pulling image %s", imgSha)
		}

		// create node
		node, _ := ns.Create(mycmder.New(ns.Name))

		ip, _, _ := node.IP()

		cfg := kube.ConfigData{
			ClusterName:          *profile,
			KubernetesVersion:    *kubeVersion,
			ControlPlaneEndpoint: ip + ":6443",
			APIBindPort:          6443,
			APIServerAddress:     *hostIP,
			Token:                "abcdef.0123456789abcdef",
			PodSubnet:            "10.244.0.0/16",
			ServiceSubnet:        "10.96.0.0/12",
			ControlPlane:         true,
			IPv6:                 false,
		}

		kCfg, _ := kube.KubeAdmCfg(cfg)

		// copy the config to the node

		if err := node.WriteFile(kube.KubeAdmCfgPath, kCfg); err != nil {
			klog.Errorf("failed to copy kubeadm config to node : %v", err)
		}

		kube.RunKubeadmInit(node, *hostIP, hostPort, *profile)
		kube.RunTaint(node)
		c, _ := kube.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container
		kube.WriteKubeConfig(c, *profile)
		kube.InstallCNI(node, "10.244.0.0/16")
	}

}

func newNodeSpec(profile string, imgSHA string, hostIP string, hostPort int32) *node.Spec {
	return &node.Spec{
		Profile:           profile,
		Name:              profile + "control-plane",
		Image:             imgSHA,
		Role:              "control-plane",
		ExtraMounts:       []cri.Mount{},
		ExtraPortMappings: []cri.PortMapping{},
		APIServerAddress:  hostIP,
		APIServerPort:     hostPort,
		IPv6:              false,
	}
}
