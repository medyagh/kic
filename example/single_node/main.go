package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/medyagh/kic/example/single_node/mycmder"
	"github.com/medyagh/kic/pkg/config/cri"
	"github.com/medyagh/kic/pkg/image"
	"github.com/medyagh/kic/pkg/kube"
	"github.com/medyagh/kic/pkg/node"
	"github.com/medyagh/kic/pkg/oci"
	"github.com/phayes/freeport"
	"k8s.io/klog"
)

func main() {
	profile := flag.String("profile", "p1", "profile name")
	delete := flag.Bool("delete", false, "to delete")
	start := flag.Bool("start", false, "to start")
	hostIP := flag.String("host-ip", "127.0.0.1", "node's ip")
	cpus := flag.String("cpu", "2", "number of cpus to dedicate to the node")
	memory := flag.String("memory", "512m", "memory")
	kubeVersion := flag.String("kubernetes-version", "v1.15.0", "kuberentes version")

	flag.Parse()
	p, err := freeport.GetFreePort()
	hostPort := int32(p)
	if err != nil {
		klog.Fatal(err)
	}

	imgSha, _ := image.NameForVersion(*kubeVersion)

	ns := &node.Spec{
		Profile:           *profile,
		Name:              *profile + "control-plane",
		Image:             imgSha,
		CPUs:              *cpus,
		Memory:            *memory,
		Role:              "control-plane",
		ExtraMounts:       []cri.Mount{},
		ExtraPortMappings: []cri.PortMapping{},
		APIServerAddress:  *hostIP,
		APIServerPort:     hostPort,
		IPv6:              false,
	}

	if *start {
		fmt.Printf("Starting on port %d\n ", hostPort)
		err := oci.PullIfNotPresent(imgSha, false, time.Minute*3)
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
		kaCfgPath := "/kic/kubeadm.conf"
		// copy the config to the node
		if err := node.WriteFile(kaCfgPath, kCfg, "644"); err != nil {
			klog.Errorf("failed to copy kubeadm config to node : %v", err)
		}

		kube.RunKubeadmInit(node, kaCfgPath, *hostIP, hostPort, *profile)
		kube.RunTaint(node)
		kube.InstallCNI(node, "10.244.0.0/16")
		c, _ := kube.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container
		// kubeconfig for end-user
		kube.WriteKubeConfig(c, *profile)
	}

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		ns.Delete()

	}

}
