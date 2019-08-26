package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	userImg := flag.String("image", "", "image to load")
	load := flag.Bool("load", false, "to load an image")

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

	cmder := mycmder.New(ns.Name)

	if *start {
		fmt.Printf("Starting on port %d\n ", hostPort)
		err := oci.PullIfNotPresent(imgSha, false, time.Minute*3)
		if err != nil {
			klog.Errorf("Error pulling image %s", imgSha)
		}

		// create node
		node, err := ns.Create(cmder)
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

		if len(*userImg) != 0 {
			loadImage(*userImg, node)
		}

		c, _ := kube.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container

		// kubeconfig for end-user
		kube.WriteKubeConfig(c, *profile)
	}

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		ns.Delete()
	}

	if *load && len(*userImg) != 0 {
		node, err := node.Find(*profile+"control-plane", cmder)
		if err != nil {
			klog.Errorf("error reading image (%s) from disk : %v", *userImg, err)
			os.Exit(1)
		}
		loadImage(*userImg, node)
	}
}

func loadImage(image string, node *node.Node) {
	_, err := oci.ImageID(image)
	if err != nil {
		klog.Errorf("error getting image not present locally %s: %v", image, err)
		os.Exit(1)
	}

	dir, err := ioutil.TempDir("", "image-tar")
	if err != nil {
		klog.Errorf("error creating temp directory")
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	imageTarPath := filepath.Join(dir, "image.tar")
	fmt.Println(imageTarPath)
	fmt.Printf("Saving image archive %s\n", image)
	err = oci.Save(image, imageTarPath)
	if err != nil {
		klog.Errorf("error saving image archive %s: %v", image, err)
		os.Exit(1)
	}

	f, err := os.Open(imageTarPath)
	if err != nil {
		klog.Errorf("error reading image (%s) from disk : %v", image, err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Printf("Loading image %s\n", image)
	err = node.LoadImageArchive(f)
	if err != nil {
		klog.Errorf("error loading (%s) into the node : %v", image, err)
		os.Exit(1)
	}
}
