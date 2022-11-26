package cmd

import (
	"fmt"
	"strings"
)

const (
	ImagePath = "%v/%v%v%v@sha256:%v"
)

// ParseImagePath parses input string to retrieve tag and digest.
// image usually a string `registry/repository-service:tag@sha256:digest`
func ParseImagePath(i string) (tag, digest string) {
	repoTagAndPrefixedDigest := strings.Split(i, "@")
	repoAndTag := strings.Split(repoTagAndPrefixedDigest[0], ":")
	if len(repoAndTag) > 1 {
		tag = repoAndTag[1]
	}
	prefixedDigest := repoTagAndPrefixedDigest[1]
	digest = strings.TrimPrefix(prefixedDigest, DigestPrefix)
	return
}

func ComposeImagePath(registry, repo, service, tag, digest string) string {
	return fmt.Sprintf(ImagePath, registry, repo, service, AppendSemicolon(tag), digest)
}
