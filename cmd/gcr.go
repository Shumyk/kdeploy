package cmd

import (
	"context"
	util "shumyk/kdeploy/cmd/util"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	ctx  = context.Background()
	auth = authn.DefaultKeychain
)

func ListRepoImages(ch chan<- *google.Tags) {
	_, err := google.NewGcloudAuthenticator()
	util.ErrorCheck(err, "GCloud authentication failed")

	gcrRepoName := FullGcrRepositoryName()
	registry := name.WithDefaultRegistry(Registry())
	repo, err := name.NewRepository(gcrRepoName, registry)
	util.ErrorCheck(err, "Failed to obtain GCR repository '"+gcrRepoName+"'")

	keychain := google.WithAuthFromKeychain(auth)
	tags, err := google.List(repo, keychain)
	util.ErrorCheck(err, "Failed to list tags for GCR repository '"+gcrRepoName+"'")

	ch <- tags
}

func ListRepos() (results []string) {
	registry, err := name.NewRegistry(Registry())
	util.ErrorCheck(err, "Failed to obtain GCR registry '"+Registry()+"'")

	authOption := remote.WithAuthFromKeychain(auth)
	repos, err := remote.Catalog(ctx, registry, authOption)
	util.ErrorCheck(err)

	return filterRepos(repos)
}

func filterRepos(reposRaw []string) (results []string) {
	for _, repoRaw := range reposRaw {
		if strings.HasPrefix(repoRaw, Repository()) {
			repo := strings.TrimPrefix(repoRaw, Repository())
			results = append(results, repo)
		}
	}
	return
}
