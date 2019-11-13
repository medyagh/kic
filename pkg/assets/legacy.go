package assets

import "io"

// LegacyCopyableFile is something that can be copied this is temproary copied from minikube to avoid circular dependecy.
type LegacyCopyableFile interface {
	io.Reader
	GetLength() int
	GetAssetName() string
	GetTargetDir() string
	GetTargetName() string
	GetPermissions() string
}
