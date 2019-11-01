package oci

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/medyagh/kic/pkg/assets"

	"github.com/pkg/errors"
)

// Copy copies a local asset into the container
func Copy(ociID string, asset assets.CopyAsset) error {
	if _, err := os.Stat(asset.AssetName); os.IsNotExist(err) {
		return errors.Wrapf(err, "error source %s does not exist", asset.AssetName)
	}

	destination := fmt.Sprintf("%s:%s", ociID, asset.TargetPath())
	cmd := exec.Command(DefaultOCI, "cp", asset.AssetName, destination)
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "error copying %s into node", asset.AssetName)
	}
	return nil
}
