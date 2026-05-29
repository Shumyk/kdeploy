package cmd

import (
	util "shumyk/kdeploy/cmd/util"
	"slices"

	garPB "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"google.golang.org/api/iterator"
)

func (g *GarGuy) ListPackages() (results []string) {
	r := &garPB.ListPackagesRequest{
		Parent:   g.conf.repositoryParentPath(),
		PageSize: 1000,
	}
	util.Debug("Listing packages using GAR: ", r.Parent)

	it := g.getClient().ListPackages(ctx, r)
	for {
		pkg, err := it.Next()
		if err == iterator.Done {
			break
		}
		util.ErrorCheck(err, "Failed to list packages in repository")

		packageName := takeLastPathPart(pkg.Name)
		results = append(results, packageToMicroserviceName(packageName))
	}

	slices.Sort(results)
	return slices.Compact(results)
}
