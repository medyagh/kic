package main

import (
	"flag"
	"fmt"
	"os"
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
	memory := flag.String("memory", "2000m", "memory")
	kubeVersion := flag.String("kubernetes-version", "v1.15.0", "kuberentes version")
	img := flag.String("image", "", "image to load")

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
		node, err := ns.Create(mycmder.New(ns.Name))
		if err != nil {
			klog.Errorf("Error Creating node %s %v", ns.Name, err)
		}

		ip, _, err := node.IP()
		if err != nil {
			klog.Errorf("Error getting node ip: %s error: %v", ip, err)
		}

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

		kCfg, err := kube.KubeAdmCfg(cfg)
		if err != nil {
			klog.Errorf("failed to generate kubeaddm  error: %v , kCfg :\n %+v", err, kCfg)
		}
		kaCfgPath := "/kic/kubeadm.conf"
		// copy the config to the node
		if err := node.WriteFile(kaCfgPath, kCfg, "644"); err != nil {
			klog.Errorf("failed to copy kubeadm config to node : %v", err)
		}

		kube.RunKubeadmInit(node, kaCfgPath, *hostIP, hostPort, *profile)
		kube.RunTaint(node)
		kube.InstallCNI(node, "10.244.0.0/16")

		if len(*img) != 0 {
			fmt.Printf("loading image %s\n", *img)
			f, err := os.Open(*img)
			if err != nil {
				klog.Errorf("error reading image (%s) from disk : %v", *img, err)
			}
			defer f.Close()
			err = node.LoadImageArchive(f)
			if err != nil {
				klog.Errorf("error loading (%s) into the node : %v", *img, err)
			}
		}

		c, _ := kube.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container
		// kubeconfig for end-user
		kube.WriteKubeConfig(c, *profile)
	}

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		ns.Delete()
	}
}
