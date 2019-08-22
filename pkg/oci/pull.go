package oci

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/medyagh/kic/pkg/exec"
)

// PullIfNotPresent pulls docker image if not present back off exponentially
func PullIfNotPresent(image string) (err error) {
	cmd := exec.Command(DefaultOCI, "inspect", "--type=image", image)
	if err := cmd.Run(); err == nil {
		return fmt.Errorf("PullIfNotPresent: image %s present locally : %v", image, err)
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
	err := exec.Command(DefaultOCI, "pull", image).Run()
	if err != nil {
		return fmt.Errorf("error pull image %s : %v", image, err)
	}
	return err
}
