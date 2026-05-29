package cmd

import (
	util "shumyk/kdeploy/cmd/util"

	gar "cloud.google.com/go/artifactregistry/apiv1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var garGuy *GarGuy

type GarGuy struct {
	conf GAR

	client *gar.Client
}

func GetGarGuy(conf GAR) *GarGuy {
	g := &GarGuy{conf: conf}
	g.buildClient()

	return g
}

func (g *GarGuy) buildClient() {
	util.Debug("Building GAR client")

	util.Debug("Finding default credentials")
	creds, err := google.FindDefaultCredentials(ctx, gar.DefaultAuthScopes()...)
	util.ErrorCheck(err, "Failed to find default credentials")

	util.Debug("Creating new GAR client")
	g.client, err = gar.NewClient(ctx, option.WithCredentials(creds))
	util.ErrorCheck(err, "Failed to create new GAR client")
}
