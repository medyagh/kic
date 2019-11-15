package image

import "fmt"

// TODO: if doesnt exist do a docker pull and try go get the sha
// NameForVersion returns the image name and sha for a kuberentes version
func NameForVersion(ver string) (string, error) {

	switch ver {
	case "v1.11.10":
		return "medyagh/kic:v1.11.10@sha256:23bb7f5e8dd2232ec829132172e87f7b9d8de65269630989e7dac1e0fe993b74", nil
	case "v1.12.8":
		return "medyagh/kic:v1.12.8@sha256:c74bc5f3efe3539f6e1ad7f11bf7c09f3091c0547cb28071f4e43067053e5898", nil
	case "v1.12.9":
		return "medyagh/kic:v1.12.9@sha256:ff82f58e18dcb22174e8eb09dae14f7edd82d91a83c7ef19e33298d0eba6a0e3", nil
	case "v1.12.10":
		return "medyagh/kic:v1.12.10@sha256:2d174bae7c20698e59791e7cca9b6db234053d1a92a009d5bb124e482540c70b", nil
	case "v1.13.6":
		return "medyagh/kic:v1.13.6@sha256:cf63e50f824fe17b90374d38d64c5964eb9fe6b3692669e1201fcf4b29af4964", nil
	case "v1.13.7":
		return "medyagh/kic:v1.13.7@sha256:1a6a5e1c7534cf3012655e99df680496df9bcf0791a304adb00617d5061233fa", nil
	case "v1.14.3":
		return "medyagh/kic:v1.14.3@sha256:cebec21f6af23d5dfa3465b88ddf4a1acb94c2c20a0a6ff8cc1c027b0a4e2cec", nil
	case "v1.15.0":
		return "medyagh/kic:v1.15.0@sha256:40d433d00a2837c8be829bd3cb0576988e377472062490bce0b18281c7f85303", nil
	case "v1.15.3":
		return "medyagh/kic:v1.15.3@sha256:f05ce52776a86c6ead806942d424de7076af3f115b0999332981a446329e6cf1", nil
	case "v1.16.1":
		return "medyagh/kic:v1.16.1@sha256:e74530d22e6a04442a97a09bdbba885ad693fcc813a0d1244da32666410d1ad1", nil
	case "v1.16.2":
		return "medyagh/kic:v1.16.2@sha256:3374a30971bf5b0011441a227fa56ef990b76125b36ca0ab8316a3c7e4f137a3", nil
	default:
		return "v1.15.0", fmt.Errorf("not supported version, using default version")
	}
}
