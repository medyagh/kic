package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/medyagh/kic/example/single_node/mycmder"
	"github.com/medyagh/kic/pkg/action"
	"github.com/medyagh/kic/pkg/assets"
	"github.com/medyagh/kic/pkg/config/cri"
	"github.com/medyagh/kic/pkg/image"
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
	copy := flag.Bool("cp", false, "to copy a file/folder into the node")
	remove := flag.Bool("rm", false, "to rm a file from the node")
	src := flag.String("src", "", "source file/folder to copy")
	dest := flag.String("dest", "", "destination to copy file/folder ")
	pause := flag.Bool("pause", false, "Pause all processes within one or more containers")

	flag.Parse()
	p, err := freeport.GetFreePort()
	hostPort := int32(p)
	if err != nil {
		klog.Fatal(err)
	}

	imgSha, err := image.NameForVersion(*kubeVersion)
	if err != nil {
		klog.Errorf("Error getting image %s", imgSha)
	}
	envs, err := getProxyEnvs()
	if err != nil {
		klog.Errorf("Error getting proxy details %v", envs)
	}
	nodeName := *profile + "-control-plane"
	ns := &node.Spec{
		Profile:           *profile,
		Name:              nodeName,
		Image:             imgSha,
		CPUs:              *cpus,
		Memory:            *memory,
		Role:              "control-plane",
		ExtraMounts:       []cri.Mount{},
		ExtraPortMappings: []cri.PortMapping{},
		APIServerAddress:  *hostIP,
		APIServerPort:     hostPort,
		IPv6:              false,
		Envs:              envs,
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

		cfg := action.ConfigData{
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

		kCfg, err := action.KubeAdmCfg(cfg)
		if err != nil {
			klog.Errorf("failed to generate kubeaddm  error: %v , kCfg :\n %+v", err, kCfg)
		}
		kaCfgPath := "/kic/kubeadm.conf"
		// copy the config to the node
		if err := node.WriteFile(kaCfgPath, kCfg, "644"); err != nil {
			klog.Errorf("failed to copy kubeadm config to node : %v", err)
		}

		_, err = action.RunKubeadmInit(node, kaCfgPath, *hostIP, hostPort, *profile)
		if err != nil {
			klog.Errorf("failed to RunKubeadmInit : %v", err)
		}

		err = action.RunTaint(node)
		if err != nil {
			klog.Errorf("failed to RunTaint : %v", err)
		}

		err = action.InstallCNI(node, "10.244.0.0/16")
		if err != nil {
			klog.Errorf("failed to InstallCNI : %v", err)
		}

		if len(*userImg) != 0 {
			loadImage(*userImg, node)
		}

		c, err := action.GenerateKubeConfig(node, *hostIP, hostPort, *profile) // generates from the /etc/ inside container
		if err != nil {
			klog.Errorf("failed to GenerateKubeConfig : %v", err)
		}

		// kubeconfig for end-user
		err = action.WriteKubeConfig(c, *profile)
		if err != nil {
			klog.Errorf("failed to WriteKubeConfig : %v", err)
		}

	}

	if *delete {
		fmt.Printf("Deleting ... %s\n", *profile)
		err = ns.Delete()
		if err != nil {
			klog.Errorf("failed to delete cluster %s : %v", *profile, err)
		}

	}

	if *load && len(*userImg) != 0 {
		node, err := node.Find(nodeName, cmder)
		if err != nil {
			klog.Errorf("error reading image (%s) from disk : %v", *userImg, err)
			os.Exit(1)
		}
		loadImage(*userImg, node)
	}

	if *copy {
		node, err := node.Find(nodeName, cmder)
		if err != nil {
			klog.Errorf("error finding node %s: %v", *userImg, err)
			os.Exit(1)
		}
		err = copyAsset(node, *src, *dest)
		if err != nil {
			klog.Errorf("error copying asset src: %s dest: %s, err: %v", *src, *dest, err)
			os.Exit(1)
		}
	}

	if *remove {
		node, err := node.Find(nodeName, cmder)
		if err != nil {
			klog.Errorf("error finding node %s: %v", *userImg, err)
			os.Exit(1)
		}
		err = node.Remove(*src)
		if err != nil {
			klog.Errorf("error removing file %s: %v", *src, *dest, err)
			os.Exit(1)
		}
	}

	if *pause {
		err := ns.Pause()
		if err != nil {
			klog.Errorf("Error pausing node %s %v", ns.Name, err)
		}
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

// getProxyEnvs returns a struct with the host environment proxy settings
// that should be passed to the nodes
func getProxyEnvs() (map[string]string, error) {
	const httpProxy = "HTTP_PROXY"
	const httpsProxy = "HTTPS_PROXY"
	const noProxy = "NO_PROXY"
	var proxyEnvs = []string{httpProxy, httpsProxy, noProxy}
	var val string
	envs := make(map[string]string)

	proxySupport := false

	for _, name := range proxyEnvs {
		val = os.Getenv(name)
		if val != "" {
			proxySupport = true
			envs[name] = val
			envs[strings.ToLower(name)] = val
		} else {
			val = os.Getenv(strings.ToLower(name))
			if val != "" {
				proxySupport = true
				envs[name] = val
				envs[strings.ToLower(name)] = val
			}
		}
	}

	// Specifically add the docker network subnets to NO_PROXY if we are using proxies
	if proxySupport {
		subnets, err := oci.GetSubnets(node.DefaultNetwork)
		if err != nil {
			return nil, err
		}
		noProxyList := strings.Join(append(subnets, envs[noProxy]), ",")
		envs[noProxy] = noProxyList
		envs[strings.ToLower(noProxy)] = noProxyList
	}

	return envs, nil
}

func copyAsset(n *node.Node, src, dest string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	targetDir := path.Dir(dest)
	targetName := path.Base(dest)

	asset := assets.CopyAsset{
		AssetName:   fileInfo.Name(),
		TargetName:  targetName,
		TargetDir:   targetDir,
		Length:      fileInfo.Size(),
		Permissions: "0777",
	}

	return n.Copy(asset)
}
