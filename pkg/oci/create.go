package oci

import (
	"github.com/medyagh/kic/pkg/config/cri"
	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
)

// CreateOpt is an option for Create
type CreateOpt func(*createOpts) *createOpts

// actual options struct
type createOpts struct {
	RunArgs       []string
	ContainerArgs []string
	Mounts        []cri.Mount
	PortMappings  []cri.PortMapping
}

// CreateContainer creates a container with "docker/podman run"
func CreateContainer(image string, opts ...CreateOpt) ([]string, error) {
	o := &createOpts{}
	for _, opt := range opts {
		o = opt(o)
	}
	// convert mounts to container run args
	runArgs := o.RunArgs
	for _, mount := range o.Mounts {
		runArgs = append(runArgs, generateMountBindings(mount)...)
	}
	for _, portMapping := range o.PortMappings {
		runArgs = append(runArgs, generatePortMappings(portMapping)...)
	}
	// construct the actual docker run argv
	args := []string{"run"}
	args = append(args, runArgs...)
	args = append(args, image)
	args = append(args, o.ContainerArgs...)
	cmd := exec.Command(DefaultOCI, args...)
	output, err := exec.CombinedOutputLines(cmd)
	// TODO : check for exist status 125 that means it alread exists, we can re-start it
	// example error:
	// $ docker run --cpus=2 --memory=2000m -d -t --privileged --security-opt seccomp=unconfined --tmpfs /tmp --tmpfs /run -v /lib/modules:/lib/modules:ro --hostname p1control-plane --name p1control-plane --label io.k8s.sigs.kic.clusterp1 --label io.k8s.sigs.kic.role=control-plane --expose 50182 --publish=127.0.0.1:50182:6443 medyagh/kic:v1.15.0@sha256:1f03b3168ffe8ab43ce170a5729e31b0d53fb3a1af88e1ad1bdf4626fad8a91c
	//		 docker: Error response from daemon: Conflict. The container name "/p1control-plane" is already in use by container "0204dcf3ca51c874b6c7dac989beae9d98dd44af53e0a17312f4d3480c1f6191". You have to remove (or rename) that container to be able to reuse that name.
	// 		 See 'docker run --help'.
	// $ echo $?
	// 125

	if err != nil {
		return output, errors.Wrapf(err, "CreateContainer %v ", args)
	}
	return output, nil
}

// WithRunArgs sets the args for docker run
// as in the args portion of `docker run args... image containerArgs...`
func WithRunArgs(args ...string) CreateOpt {
	return func(r *createOpts) *createOpts {
		r.RunArgs = args
		return r
	}
}

// WithMounts sets the container mounts
func WithMounts(mounts []cri.Mount) CreateOpt {
	return func(r *createOpts) *createOpts {
		r.Mounts = mounts
		return r
	}
}

// WithPortMappings sets the container port mappings to the host
func WithPortMappings(portMappings []cri.PortMapping) CreateOpt {
	return func(r *createOpts) *createOpts {
		r.PortMappings = portMappings
		return r
	}
}
