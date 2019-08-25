package node

import (
	"github.com/medyagh/kic/pkg/runner"
)

// Find finds a node
func Find(name string, cmder runner.Cmder) (*Node, error) {
	// TODO: check node exists

	return &Node{
		name:  name,
		cache: &nodeCache{},
		cmder: cmder,
	}, nil
}
