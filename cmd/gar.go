package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"
	"strings"

	gar "cloud.google.com/go/artifactregistry/apiv1"
	garPB "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func ListRepoImagesGAR(ch chan<- model.ImageOptions) {
	util.Debug("Listing images in repository using GAR")

	util.Debug("Finding default credentials")
	creds, err := google.FindDefaultCredentials(ctx, gar.DefaultAuthScopes()...)
	util.ErrorCheck(err, "Failed to find default credentials")

	util.Debug("Creating new GAR client")
	client, err := gar.NewClient(ctx, option.WithCredentials(creds))
	util.ErrorCheck(err, "Failed to create new GAR client")

	defer client.Close()

	req := &garPB.ListVersionsRequest{
		Parent:   "projects/" + config.GAR.Project + "/locations/" + config.GAR.Location + "/repositories/" + config.GAR.Repository + "/packages/" + GarPackageName(),
		PageSize: 100,
		OrderBy:  "create_time desc",
		View:     garPB.VersionView_FULL,
	}

	util.Debug("Listing images in repository: ", req.Parent)
	it := client.ListVersions(ctx, req)
	util.Debug("Retrieved images: ", it.PageInfo().MaxSize)

	var images model.ImageOptions
	for {
		if len(images) > int(req.PageSize) {
			break
		}

		util.Debug("Reading next image")
		image, err := it.Next()

		if err != nil {
			util.Laugh(err, "Failed to list images in repository")
			break
		}

		util.Debug("Image: ", image.Name, image.RelatedTags)

		// tags format is:
		// 		projects/{project}/locations/{location}/repositories/{repo}/packages/{package}/tags/{tag}
		// so, we need to take only {tag} part
		tags := util.SliceMapping(image.RelatedTags, func(t *garPB.Tag) string {
			parts := strings.Split(t.Name, "/")
			return parts[len(parts)-1]
		})
		images = append(
			images,
			model.ImageOption{
				Created: image.UpdateTime.AsTime(),
				Tags:    tags,
				Digest:  strings.Split(image.Name, ":")[1],
			},
		)
	}

	ch <- images
}
