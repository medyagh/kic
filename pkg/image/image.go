package image

import "fmt"

// TODO: if doesnt exist do a docker pull and try go get the sha
// NameForVersion returns the image name and sha for a kuberentes version
func NameForVersion(ver string) (string, error) {

	switch ver {
	case "v1.11.10":
		return "medyagh/kic:v1.11.10@sha256:6cfec38db507b9b35b8779010640fe7f59d9e27b111cb9079fb42d801116d090", nil
	case "v1.12.8":
		return "medyagh/kic:v1.12.8@sha256:fb4827ffa316b421985e6c279048195b6d207feed25d0b0d20d3b0df0d78206c", nil
	case "v1.12.9":
		return "medyagh/kic:v1.12.9@sha256:518265de75580eda358f91789e720c1a162f4ce131ee7b605d26836a06967f20", nil
	case "v1.12.10":
		return "medyagh/kic:v1.12.10@sha256:f565c866813e4184d3454ff0d576b4dcdeddf21b9b9ee5da3fcd88099967ead9", nil
	case "v1.13.6":
		return "medyagh/kic:v1.13.6@sha256:0fd2095ccdd145acbc05922485c99703e7378e3242f9fa78197e71368ee464cb", nil
	case "v1.13.7":
		return "medyagh/kic:v1.13.7@sha256:f10a8c4fde2a3849fb5954294457c8cd7a4baae075cbcdf115ca286430acb4eb", nil
	case "v1.14.3":
		return "medyagh/kic:v1.14.3@sha256:8f7be744af5776c47247202fc08bd4eaf7173178ce26038b40f234f682d46192", nil
	case "v1.15.0":
		return "medyagh/kic:v1.15.0@sha256:103f0592f6073d4d4ea1167e7e487d4555a25d140d297fc2d2a3fb2ec3ecd4d7", nil
	case "v1.15.3":
		return "medyagh/kic:v1.15.3@sha256:63358cd2952dc63d37c72b118a745f8080a72f7a8f3928ffa721358eb85cd14b", nil
	case "v1.16.1":
		return "medyagh/kic:v1.16.1@sha256:ca01b2104c94422fdbf55833d6fce2192b9d15fd064662cbfb0e542bee8c9898", nil
	case "v1.16.2":
		return "medyagh/kic:v1.16.2@sha256:2335bde84ead2ff9543c6729aa5edf16f77cc66c9e4feff69420eea697399391", nil
	default:
		return "v1.15.0", fmt.Errorf("not supported version, using default version")
	}
}
