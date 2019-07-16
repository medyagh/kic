package kube

import (
	"fmt"
	"strings"

	"github.com/medyagh/kic/pkg/config"
	"github.com/medyagh/kic/pkg/config/kustomize"
)

func GetMagicConfig(nodeIP string, profile string, kubeVersion string) (string, error) {
	clusterCfg := &config.Cluster{}
	cd := ConfigData{
		ClusterName:          profile,
		KubernetesVersion:    kubeVersion,
		ControlPlaneEndpoint: nodeIP + ":6443",
		APIBindPort:          6443,
		APIServerAddress:     "127.0.0.1",
		Token:                "abcdef.0123456789abcdef",
		PodSubnet:            "10.244.0.0/16",
		ServiceSubnet:        "10.96.0.0/12",
		ControlPlane:         true,
		IPv6:                 false,
	}

	config, err := Config(cd)
	if err != nil {
		return "", err
	}
	// fix all the patches to have name metadata matching the generated config
	patches, jsonPatches := setPatchNames(
		allPatchesFromConfig(clusterCfg),
	)
	// apply patches
	patched, err := kustomize.Build([]string{config}, patches, jsonPatches)
	if err != nil {
		return "", err
	}
	return removeMetadata(patched), nil

	// return writeKubeadmConfig(clusterConfig, configData, node)
}

// trims out the metadata.name we put in the config for kustomize matching,
// kubeadm will complain about this otherwise
func removeMetadata(kustomized string) string {
	return strings.Replace(
		kustomized,
		`metadata:
  name: config
`,
		"",
		-1,
	)
}

func allPatchesFromConfig(cfg *config.Cluster) (patches []string, jsonPatches []kustomize.PatchJSON6902) {
	return cfg.KubeadmConfigPatches, cfg.KubeadmConfigPatchesJSON6902
}

// setPatchNames sets the targeted object name on every patch to be the fixed
// name we use when generating config objects (we have one of each type, all of
// which have the same fixed name)
func setPatchNames(patches []string, jsonPatches []kustomize.PatchJSON6902) ([]string, []kustomize.PatchJSON6902) {
	fixedPatches := make([]string, len(patches))
	fixedJSONPatches := make([]kustomize.PatchJSON6902, len(jsonPatches))
	for i, patch := range patches {
		// insert the generated name metadata
		fixedPatches[i] = fmt.Sprintf("metadata:\nname: %s\n%s", ObjectName, patch)
	}
	for i, patch := range jsonPatches {
		// insert the generated name metadata
		patch.Name = ObjectName
		fixedJSONPatches[i] = patch
	}
	return fixedPatches, fixedJSONPatches
}
