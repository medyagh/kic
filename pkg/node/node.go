package node

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/medyagh/kic/pkg/oci"
	"github.com/medyagh/kic/pkg/runner"

	"github.com/pkg/errors"
)

// Node represents a handle to a kic node
// This struct must be created by one of: CreateControlPlane
type Node struct {
	// must be one of docker container ID or name
	name string
	// cached node info etc.
	cache *nodeCache
	cmder runner.Cmder
}

// WriteFile writes content to dest on the node
func (n *Node) WriteFile(dest, content string, perm string) error {
	// create destination directory
	cmd := n.Command("mkdir", "-p", filepath.Dir(dest))
	_, err := runner.RunLoggingOutputOnFail(cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dest)
	}

	err = n.Command("cp", "/dev/stdin", dest).SetStdin(strings.NewReader(content)).Run()
	if err != nil {
		return errors.Wrapf(err, "failed to run: cp /dev/stdin %s", dest)
	}
	err = n.Command("chmod", perm, dest).Run()
	return errors.Wrapf(err, "failed to run: chmod %s %s", perm, dest)
}

// IP returns the IP address of the node
func (n *Node) IP() (ipv4 string, ipv6 string, err error) {
	// use the cached version first
	cachedIPv4, cachedIPv6 := n.cache.IP()
	if cachedIPv4 != "" && cachedIPv6 != "" {
		return cachedIPv4, cachedIPv6, nil
	}
	// retrieve the IP address of the node using docker inspect
	lines, err := oci.Inspect(n.name, "{{range .NetworkSettings.Networks}}{{.IPAddress}},{{.GlobalIPv6Address}}{{end}}")
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get container details")
	}
	if len(lines) != 1 {
		return "", "", errors.Errorf("file should only be one line, got %d lines", len(lines))
	}
	ips := strings.Split(lines[0], ",")
	if len(ips) != 2 {
		return "", "", errors.Errorf("container addresses should have 2 values, got %d values", len(ips))
	}
	n.cache.set(func(cache *nodeCache) {
		cache.ipv4 = ips[0]
		cache.ipv6 = ips[1]
	})
	return ips[0], ips[1], nil
}

// LoadImageArchive loads an image form archive into node
func (n *Node) LoadImageArchive(image io.Reader) error {
	cmd := n.Command(
		"ctr", "--namespace=k8s.io", "images", "import", "-",
	)
	cmd.SetStdin(image)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to load image")
	}
	return nil
}

// Command returns a new runner.Cmd that will run on the node
func (n *Node) Command(command string, args ...string) runner.Cmd {
	return n.cmder.Command(command, args...)
}
