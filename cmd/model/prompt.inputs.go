package model

import (
	"github.com/google/go-containerregistry/pkg/v1/google"
	util "shumyk/kdeploy/cmd/util"
)

type PromptInputs interface {
	ImageOptions() ImageOptions
}

type PreviousImages []PreviousImage

func (i PreviousImages) ImageOptions() ImageOptions {
	return util.SliceMapping(i, PreviousImage.ImageOption)
}

type Manifests map[string]google.ManifestInfo

func (m Manifests) ImageOptions() ImageOptions {
	return util.MapToSliceMapping(m, ManifestInfoToImageOption)
}
