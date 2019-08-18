package image

// NameForVersion returns the image name and sha for a kuberentes version
func NameForVersion(ver string) (string, error) {

	switch ver {
	case "v1.15.0":
		return "medyagh/kic:v1.15.0@sha256:ef287ad027aa9b029edf74e4288e7673654fc3c93ffe4779e6a4005d9d4b030e", nil
	case "v1.13.0":
		return "kindest/node:v1.13.0@sha256:b4d092fd2b507843dd096fe6c85d06a27a0cbd740a0b32a880fe61aba24bb478", nil
	default:
		return "v1.15.0", nil
	}
}
