package model

import (
	"sort"
	"time"

	util "shumyk/kdeploy/cmd/util"
)

type ImageOption struct {
	Created time.Time
	Tags    []string
	Digest  string
}

func (o ImageOption) String() string {
	return util.FormatImageOption(o.Created, o.Digest, o.Tags...)
}

type ImageOptions []ImageOption

func (o ImageOptions) ImageOptions() ImageOptions {
	return o
}

func (o ImageOptions) Stringify() []string {
	return util.SliceMapping(o, ImageOption.String)
}

func (o ImageOptions) Sorted() ImageOptions {
	sort.SliceStable(o, sortByCreated(o))
	return o
}

func sortByCreated(o ImageOptions) func(i, j int) bool {
	return func(i, j int) bool {
		return o[i].Created.After(o[j].Created)
	}
}
