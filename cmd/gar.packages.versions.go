package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"

	garPB "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"

	"google.golang.org/api/iterator"
)

func (g GarGuy) ListPackageVersions(ch chan<- model.ImageOptions) {
	req := &garPB.ListVersionsRequest{
		Parent:   g.conf.packageParentPath(GarPackageName()),
		PageSize: 500,
		OrderBy:  "create_time desc",
		View:     garPB.VersionView_FULL,
	}
	util.Debug("Listing versions of package using GAR ", req.Parent)

	it := g.getClient().ListVersions(ctx, req)
	util.Debug("Retrieved images: ", it.PageInfo().MaxSize)

	var images model.ImageOptions
	for {
		if len(images) >= int(req.PageSize) {
			break
		}

		image, err := it.Next()
		if err == iterator.Done {
			util.Debug(err, "Iterator is done")
			break
		}
		util.ErrorCheck(err, "Failed to list package versions")

		images = append(
			images,
			model.ImageOption{
				Created: image.UpdateTime.AsTime(),
				Tags:    takeTagNames(image.RelatedTags),
				Digest:  takeDigest(image.Name),
			},
		)
	}

	ch <- images
}

func takeTagNames(tags []*garPB.Tag) []string {
	mapper := func(t *garPB.Tag) string {
		return takeLastPathPart(t.Name)
	}
	return util.SliceMapping(tags, mapper)
}
