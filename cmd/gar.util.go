package cmd

import "strings"

// paths are formatted like:
//
//	projects/{project}/locations/{location}/repositories/{repo}/packages/{package}
//	projects/{project}/locations/{location}/repositories/{repo}/packages/{package}/tags/{tag}
//
// so, this method helps to take only last part - {package} or {tag}
func takeLastPathPart(resourceName string) string {
	parts := strings.Split(resourceName, "/")
	return parts[len(parts)-1]
}

func packageToMicroserviceName(packageName string) string {
	for serviceName, mapping := range config.Mappings {
		if mapping.GAR == packageName {
			return serviceName
		}
	}

	result, _ := strings.CutPrefix(packageName, config.GAR.PackagePrefix)
	return result
}

// image name format:
//
//	projects/{project}/locations/{location}/repositories/{repo}/packages/{package}/versions/sha256:{digest}
//
// so, we need only digest after `sha256:`
func takeDigest(i string) string {
	return strings.Split(i, ":")[1]
}
