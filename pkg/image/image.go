package image

import "fmt"

// TODO: if doesnt exist do a docker pull and try go get the sha
// NameForVersion returns the image name and sha for a kuberentes version
func NameForVersion(ver string) (string, error) {

	switch ver {
	case "v1.11.10":
		return "medyagh/kic:v1.11.10@sha256:dae9633404059f3b1debc0ca016e1a2efcfd3123793264422a873f074beb3ebc", nil
	case "v1.12.8":
		return "medyagh/kic:v1.12.8@sha256:a21a4387d660c6abbc0857174bdec365ec7f6927a54e166818e8adecfa037d6c", nil
	case "v1.12.9":
		return "medyagh/kic:v1.12.9@sha256:07142d61b988aec8706dc71d34344f31e1cfac92410830a5d34b3e9012a02757", nil
	case "v1.12.10":
		return "medyagh/kic:v1.12.10@sha256:1da46fe83a3bc918288c7ed9133105a8dea95d141a6de018b4b13f1cad6808da", nil
	case "v1.13.6":
		return "medyagh/kic:v1.13.6@sha256:294f24dc8bee373d4343cdf29e4223b134dda977d830d6140f6d4db55eda39aa", nil
	case "v1.13.7":
		return "medyagh/kic:v1.13.7@sha256:6012d101f19bd502981591027ea9dfc2e660c039b5b4e5e6ad95a5af41db47a0", nil
	case "v1.14.3":
		return "medyagh/kic:v1.14.3@sha256:533dc16fd161244cc053b88d8091878efae5437b141f51d9b3bf7069c05f4725", nil
	case "v1.15.0":
		return "medyagh/kic:v1.15.0@sha256:8f9e16c26f65e23ed6637c67be4e4cec287ce679f9c79f3eae1a09d8eeb7c853", nil
	case "v1.15.3":
		return "medyagh/kic:v1.15.3@sha256:8a3694870e5eb2123e1ad189694702ad64876fa22f5f0b3bea042129bd296f93", nil
	case "v1.16.1":
		return "medyagh/kic:v1.16.1@sha256:6bbf560b7e9a7acb241e5f2e7c063ed6c9c87387f4b1657da73d98349c315049", nil
	case "v1.16.2":
		return "medyagh/kic:v1.16.2@sha256:9880a59ee94573cf257a5118ffc8dd4f3827229fdb41c59d6b6f785edc215f54", nil
	default:
		return "v1.15.0", fmt.Errorf("not supported version, using default version")
	}
}
