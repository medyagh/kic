package config

import (
	"github.com/medyagh/kic/pkg/config/kustomize"
	"github.com/medyagh/kic/pkg/node"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cluster contains cluster configuration
type Cluster struct {
	// TypeMeta representing the type of the object and its API schema version.
	metav1.TypeMeta `json:",inline"`

	// Nodes contains the list of nodes defined in the `kic` Cluster
	// If unset this will default to a single control-plane node
	// Note that if more than one control plane is specified, an external
	// control plane load balancer will be provisioned implicitly
	Nodes []node.Node `json:"nodes"`

	/* Advanced fields */

	// Networking contains cluster wide network settings
	Networking Networking `json:"networking"`

	// KubeadmConfigPatches are applied to the generated kubeadm config as
	// strategic merge patches to `kustomize build` internally
	// https://github.com/kubernetes/community/blob/master/contributors/devel/strategic-merge-patch.md
	// This should be an inline yaml blob-string
	KubeadmConfigPatches []string `json:"kubeadmConfigPatches,omitempty"`

	// KubeadmConfigPatchesJSON6902 are applied to the generated kubeadm config
	// as patchesJson6902 to `kustomize build`
	KubeadmConfigPatchesJSON6902 []kustomize.PatchJSON6902 `json:"kubeadmConfigPatchesJson6902,omitempty"`
}

// Networking contains cluster wide network settings
type Networking struct {
	// IPFamily is the network cluster model, currently it can be ipv4 or ipv6
	IPFamily ClusterIPFamily `json:"ipFamily,omitempty"`
	// APIServerPort is the listen port on the host for the Kubernetes API Server
	// Defaults to a random port on the host
	APIServerPort int32 `json:"apiServerPort,omitempty"`
	// APIServerAddress is the listen address on the host for the Kubernetes
	// API Server. This should be an IP address.
	//
	// Defaults to 127.0.0.1
	APIServerAddress string `json:"apiServerAddress,omitempty"`
	// PodSubnet is the CIDR used for pod IPs
	// kicd will select a default if unspecified
	PodSubnet string `json:"podSubnet,omitempty"`
	// ServiceSubnet is the CIDR used for services VIPs
	// kinc will select a default if unspecified for IPv6
	ServiceSubnet string `json:"serviceSubnet,omitempty"`
	// If DisableDefaultCNI is true, kic will not install the default CNI setup.
	// Instead the user should install their own CNI after creating the cluster.
	DisableDefaultCNI bool `json:"disableDefaultCNI,omitempty"`
}

// ClusterIPFamily defines cluster network IP family
type ClusterIPFamily string

const (
	// IPv4Family sets ClusterIPFamily to ipv4
	IPv4Family ClusterIPFamily = "ipv4"
	// IPv6Family sets ClusterIPFamily to ipv6
	IPv6Family ClusterIPFamily = "ipv6"
)
