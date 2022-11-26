package model

import (
	util "shumyk/kdeploy/cmd/util"
	"strings"
)

type SelectedImage struct {
	Tags   []string
	Digest string
}

func (i SelectedImage) Tag() string {
	if len(i.Tags) != 0 {
		return i.Tags[0]
	}
	return ""
}

// ParseSelectedImage parses selected image string to corresponding struct
func ParseSelectedImage(value string) (i SelectedImage) {
	// contents of array:
	//		0: created/deployed time, not needed anymore
	//		1: digest
	//		2: tags, optional
	selectedImageData := strings.Split(value, util.Divider)
	i.Digest = selectedImageData[1]
	i.Tags = strings.Split(selectedImageData[2], util.Delimiter)
	return
}
