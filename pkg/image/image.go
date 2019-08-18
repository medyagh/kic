package image

import "fmt"

// NameForVersion returns the image name and sha for a kuberentes version
func NameForVersion(ver string) (string, error) {

	switch ver {
	case "v1.15.0":
		return "medyagh/kic:v1.15.0@sha256:1f03b3168ffe8ab43ce170a5729e31b0d53fb3a1af88e1ad1bdf4626fad8a91c", nil
	case "v1.14.3":
		return "medyagh/kic:v1.13.0@sha256:438cbaa4606e86db814d2392a6bdd449f556fe477ccfdec27964d2bbcec3e368", nil
	case "v1.13.7":
		return "medyagh/kic:v1.13.0@sha256:e07125761e5592def87e1c2945306203cac28d473cdf9b48e4c67e889c06829f", nil
	case "v1.12.9":
		return "medyagh/kic:v1.13.0@sha256:311cadf11241cf29e168cd35c8b42cce10321480f194b602395965642ec0319a", nil
	case "v1.11.10":
		return "medyagh/kic:v1.13.0@sha256:e4d808a30be9d87cd2266819bed6377a83c0835830b02043291dc359c002b1e4", nil
	default:
		return "v1.15.0", fmt.Errorf("Not supported version, using default version")
	}
}
