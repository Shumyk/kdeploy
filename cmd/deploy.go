package cmd

import (
	. "shumyk/kdeploy/cmd/model"
	prompt "shumyk/kdeploy/cmd/prompt"
	. "shumyk/kdeploy/cmd/util"

	"github.com/google/go-containerregistry/pkg/v1/google"
)

type ImageSelecter func(<-chan bool) SelectedImage

func deployTemplate(selectImage ImageSelecter) {
	clientConfig := CreateClientConfigFromKubeConfig()
	go LoadMetadata(clientConfig)

	clientSetCreatedChannel := make(chan bool)
	go ClientSet(clientConfig, clientSetCreatedChannel)

	selectedImage := selectImage(clientSetCreatedChannel)

	SetImage(&selectedImage)
}

func DeployNew() {
	deployTemplate(newImageSelecter)
}

func newImageSelecter(clientSetCreated <-chan bool) SelectedImage {
	images := make(chan *google.Tags)
	go ListRepoImages(images)

	<-clientSetCreated
	tag, digest := GetImage()
	defer SaveDeployedImage(tag, digest)
	PrintImageInfo(HeaderCurrentImage, tag, digest)

	var manifests Manifests = (<-images).Manifests
	return prompt.ImageSelect(manifests)
}

func DeployPrevious(images PreviousImages) {
	deployTemplate(previousImageSelecter(images))
}

func previousImageSelecter(images PreviousImages) ImageSelecter {
	return func(clientSetCreated <-chan bool) SelectedImage {
		selected := prompt.ImageSelect(images)
		// wait for client set to be created, as next line outside is setting image
		<-clientSetCreated
		return selected
	}
}
