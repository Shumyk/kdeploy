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

func (opts ImageOptions) Stringify() []string {
	return util.SliceMapping(opts, ImageOption.String)
}

func (opts ImageOptions) Sorted() ImageOptions {
	sort.SliceStable(opts, sortByCreated(opts))
	return opts
}

func sortByCreated(options ImageOptions) func(i, j int) bool {
	return func(i, j int) bool {
		return options[i].Created.After(options[j].Created)
	}
}
