package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	prompt "shumyk/kdeploy/cmd/prompt"
	util "shumyk/kdeploy/cmd/util"
)

type ImageSelecter func(<-chan bool) model.SelectedImage

func deployTemplate(selectImage ImageSelecter) {
	util.Debug("Creating client config from K8S config")
	clientConfig := CreateClientConfigFromKubeConfig()
	go LoadMetadata(clientConfig)

	util.Debug("Creating client set")
	clientSetCreatedChannel := make(chan bool)
	go ClientSet(clientConfig, clientSetCreatedChannel)

	util.Debug("Selecting image")
	selectedImage := selectImage(clientSetCreatedChannel)
	util.Debug("Selected Image: ", selectedImage)

	SetImage(&selectedImage)
}

func DeployNew() {
	deployTemplate(newImageSelecter)
}

func newImageSelecter(clientSetCreated <-chan bool) model.SelectedImage {
	images := make(chan model.ImageOptions)
	go garGuy.ListPackageVersions(images)

	<-clientSetCreated
	tag, digest := GetImage()
	defer SaveDeployedImage(tag, digest)
	util.PrintImageInfo(util.HeaderCurrentImage, tag, digest)

	return prompt.ImageSelect(<-images)
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
