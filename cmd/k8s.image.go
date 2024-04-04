package cmd

import (
	"encoding/json"
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"

	"k8s.io/apimachinery/pkg/types"
	confApps "k8s.io/client-go/applyconfigurations/apps/v1"
	core "k8s.io/client-go/applyconfigurations/core/v1"
)

func GetImage() (tag, digest string) {
	var response model.K8sResourceAgnosticResponse
	err := clientSet.AppsV1().RESTClient().
		Get().
		Namespace(k8sNamespace).
		Resource(k8sResourceType).
		Name(k8sResourceFullName).
		Do(ctx).
		Into(&response)
	util.ErrorCheck(err, "GET image failed")

	imagePath := response.Spec.Template.Spec.Containers[0].Image
	return util.ParseImagePath(imagePath)
}

func SetImage(image *model.SelectedImage) {
	newImage := util.ComposeImagePath(Registry(), Repository(), GcrRepositoryName(), image.Tag(), image.Digest)
	util.Debug("Setting new image: ", newImage)

	imagePatch := composeImagePatch(newImage)
	data, err := json.Marshal(imagePatch)
	util.ErrorCheck(err, "Marshalling image patch failed")

	updateError := clientSet.AppsV1().RESTClient().
		Patch(types.StrategicMergePatchType).
		Namespace(k8sNamespace).
		Resource(k8sResourceType).
		Name(k8sResourceFullName).
		Body(data).
		Do(ctx).
		Error()

	util.ErrorCheck(updateError, "PATCH image failed")
	util.PrintImageInfo(util.HeaderDeployedImage, image.Tags[0], image.Digest)
}

// composeImagePatch composes a resource apply configuration to patch only the image of a Kubernetes resource.
// It uses DeploymentApplyConfiguration, but it's actually resource agnostic as it only patches the image,
// which is located in the same place among resources.
func composeImagePatch(newImage string) confApps.DeploymentApplyConfiguration {
	containerName := ContainerName()
	container := core.ContainerApplyConfiguration{Image: &newImage, Name: &containerName}
	podSpec := core.PodSpecApplyConfiguration{Containers: []core.ContainerApplyConfiguration{container}}
	templateSpec := core.PodTemplateSpecApplyConfiguration{Spec: &podSpec}
	resourceSpec := confApps.DeploymentSpecApplyConfiguration{Template: &templateSpec}
	return confApps.DeploymentApplyConfiguration{Spec: &resourceSpec}
}
