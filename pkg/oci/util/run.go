package util

import (
	"github.com/medyagh/kic/pkg/exec"
	"k8s.io/klog"
)

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
	cmd := exec.Command("docker", args...)
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

// WithContainerArgs sets the args to the container
// as in the containerArgs portion of `docker run args... image containerArgs...`
// NOTE: this is only the args portion before the image
func WithContainerArgs(args ...string) RunOpt {
	return func(r *runOpts) *runOpts {
		r.ContainerArgs = args
		return r
	}
}
