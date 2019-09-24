package assets

import (
	"path"
)

// CopyAsset is something that can be copied
type CopyAsset struct {
	Length      int64
	AssetName   string
	TargetDir   string
	TargetName  string
	Permissions string
}

// TargetPath returns asset TargetDir/TargetName joined
func (a CopyAsset) TargetPath() string {
	return path.Join(a.TargetDir, a.TargetName)
}
