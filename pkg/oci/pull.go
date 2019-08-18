package oci

import (
	"time"

	"github.com/cenkalti/backoff"
	"k8s.io/klog"

	"github.com/medyagh/kic/pkg/exec"
)

// PullIfNotPresent pulls docker image if not present back off exponentially
func PullIfNotPresent(image string) (err error) {
	cmd := exec.Command(DefaultOCI, "inspect", "--type=image", image)
	if err := cmd.Run(); err == nil {
		klog.Infof("Image: %s present locally", image)
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute

	f := func() error {
		return pull(image)
	}

	err = backoff.Retry(f, b)

	return err
}

// Pull pulls an image, retrying up to retries times
func pull(image string) error {
	klog.Infof("Trying to pulling image: %s ...", image)
	err := exec.Command(DefaultOCI, "pull", image).Run()
	if err != nil {
		klog.Errorf("Temproary error : %v Trying again to pull image: %s ...", err, image)
	}
	return err
}
