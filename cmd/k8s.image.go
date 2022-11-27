package cmd

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/types"
	confApps "k8s.io/client-go/applyconfigurations/apps/v1"
	core "k8s.io/client-go/applyconfigurations/core/v1"
	. "shumyk/kdeploy/cmd/model"
	. "shumyk/kdeploy/cmd/util"
)

func GetImage() (tag, digest string) {
	var response K8sResourceAgnosticResponse
	err := clientSet.AppsV1().RESTClient().
		Get().
		Namespace(namespace).
		Resource(k8sResource).
		Name(k8sResourceName).
		Do(ctx).
		Into(&response)
	ErrorCheck(err, "GET image failed")

	imagePath := response.Spec.Template.Spec.Containers[0].Image
	return ParseImagePath(imagePath)
}

func SetImage(image *SelectedImage) {
	newImage := ComposeImagePath(Registry(), Repository(), microservice, image.Tag(), image.Digest)
	imagePatch := composeImagePatch(newImage)
	data, err := json.Marshal(imagePatch)
	ErrorCheck(err, "Unmarshalling image patch failed")

	updateError := clientSet.AppsV1().RESTClient().
		Patch(types.StrategicMergePatchType).
		Namespace(namespace).
		Resource(k8sResource).
		Name(k8sResourceName).
		Body(data).
		Do(ctx).
		Error()

	ErrorCheck(updateError, "PATCH image failed")
	PrintImageInfo(HeaderDeployedImage, image.Tags[0], image.Digest)
}

// composeImagePatch composes resource apply configuration to patch only image.
// DeploymentApplyConfiguration is used, but it's actually resource agnostic as we patch only image,
// which is located under same place among resources.
func composeImagePatch(newImage string) confApps.DeploymentApplyConfiguration {
	container := core.ContainerApplyConfiguration{Image: &newImage, Name: &microservice}
	podSpec := core.PodSpecApplyConfiguration{Containers: []core.ContainerApplyConfiguration{container}}
	templateSpec := core.PodTemplateSpecApplyConfiguration{Spec: &podSpec}
	resourceSpec := confApps.DeploymentSpecApplyConfiguration{Template: &templateSpec}
	return confApps.DeploymentApplyConfiguration{Spec: &resourceSpec}
}
