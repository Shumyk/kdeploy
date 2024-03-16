package model

import "time"

type PreviousImage struct {
	Tag      string
	Digest   string
	Deployed time.Time
}

func (p PreviousImage) ImageOption() ImageOption {
	return ImageOption{
		Created: p.Deployed,
		Tags:    []string{p.Tag},
		Digest:  p.Digest,
	}
}

func PreviousImageOf(tag, digest string) PreviousImage {
	return PreviousImage{
		Tag:      tag,
		Digest:   digest,
		Deployed: time.Now(),
	}
}
