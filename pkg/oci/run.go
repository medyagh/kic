package oci

import (
	"github.com/medyagh/kic/pkg/exec"
	"github.com/medyagh/kic/pkg/node/cri"
	"k8s.io/klog"
)

// RunOpt is an option for Run
type RunOpt func(*runOpts) *runOpts

// actual options struct
type runOpts struct {
	RunArgs       []string
	ContainerArgs []string
	Mounts        []cri.Mount
	PortMappings  []cri.PortMapping
}

// Run creates a container with "docker run", with some error handling
func Run(image string, opts ...RunOpt) error {
	o := &runOpts{}
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
	if err != nil {
		// log error output if there was any
		for _, line := range output {
			klog.Error(line)
		}
		return err
	}
	return nil
}

// WithRunArgs sets the args for docker run
// as in the args portion of `docker run args... image containerArgs...`
func WithRunArgs(args ...string) RunOpt {
	return func(r *runOpts) *runOpts {
		r.RunArgs = args
		return r
	}
}

// WithMounts sets the container mounts
func WithMounts(mounts []cri.Mount) RunOpt {
	return func(r *runOpts) *runOpts {
		r.Mounts = mounts
		return r
	}
}

// WithPortMappings sets the container port mappings to the host
func WithPortMappings(portMappings []cri.PortMapping) RunOpt {
	return func(r *runOpts) *runOpts {
		r.PortMappings = portMappings
		return r
	}
}
