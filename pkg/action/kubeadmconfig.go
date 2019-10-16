package action

import (
	"fmt"
	"strings"

	"github.com/medyagh/kic/pkg/cluster"
	"github.com/medyagh/kic/pkg/config/kustomize"
)

// KubeAdmCfg returns the kubeadm config
func KubeAdmCfg(cd ConfigData) (string, error) {
	clusterCfg := &cluster.Config{}
	config, err := templateExec(cd)
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

func allPatchesFromConfig(cfg *cluster.Config) (patches []string, jsonPatches []kustomize.PatchJSON6902) {
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
