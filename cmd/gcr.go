package cmd

import (
	"context"
	. "shumyk/kdeploy/cmd/util"
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
	ErrorCheck(err, "GCloud authentication failed")

	registry := name.WithDefaultRegistry(Registry())
	repo, err := name.NewRepository(BuildRepository(microservice), registry)
	ErrorCheck(err, "Obtaining new repository failed")

	keychain := google.WithAuthFromKeychain(auth)
	tags, err := google.List(repo, keychain)
	ErrorCheck(err, "Listing tags failed")

	ch <- tags
}

func ListRepos() (results []string) {
	registry, err := name.NewRegistry(Registry())
	ErrorCheck(err, "Obtaining new registry failed")

	authOption := remote.WithAuthFromKeychain(auth)
	repos, err := remote.Catalog(ctx, registry, authOption)
	ErrorCheck(err)

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
