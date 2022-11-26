package model

import "github.com/google/go-containerregistry/pkg/v1/google"

func ManifestInfoToImageOption(digest string, manifest google.ManifestInfo) ImageOption {
	return ImageOption{
		Created: manifest.Created,
		Tags:    manifest.Tags,
		Digest:  digest,
	}
}
