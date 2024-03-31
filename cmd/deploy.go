package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	prompt "shumyk/kdeploy/cmd/prompt"
	util "shumyk/kdeploy/cmd/util"

	"github.com/google/go-containerregistry/pkg/v1/google"
)

type ImageSelecter func(<-chan bool) model.SelectedImage

func deployTemplate(selectImage ImageSelecter) {
	clientConfig := CreateClientConfigFromKubeConfig()
	go LoadMetadata(clientConfig)

	clientSetCreatedChannel := make(chan bool)
	go ClientSet(clientConfig, clientSetCreatedChannel)

	selectedImage := selectImage(clientSetCreatedChannel)
	util.Debug("Selected Image: ", selectedImage)

	SetImage(&selectedImage)
}

func DeployNew() {
	deployTemplate(newImageSelecter)
}

func newImageSelecter(clientSetCreated <-chan bool) model.SelectedImage {
	images := make(chan *google.Tags)
	go ListRepoImages(images)

	<-clientSetCreated
	tag, digest := GetImage()
	defer SaveDeployedImage(tag, digest)
	util.PrintImageInfo(util.HeaderCurrentImage, tag, digest)

	var manifests model.Manifests = (<-images).Manifests
	return prompt.ImageSelect(manifests)
}

func DeployPrevious(images model.PreviousImages) {
	deployTemplate(previousImageSelecter(images))
}

func previousImageSelecter(images model.PreviousImages) ImageSelecter {
	return func(clientSetCreated <-chan bool) model.SelectedImage {
		selected := prompt.ImageSelect(images)
		// wait for client set to be created, as next line outside is setting image
		<-clientSetCreated
		return selected
	}
}
