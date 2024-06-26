package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"

	"github.com/spf13/viper"
)

func SetConfig(key string, value any) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func SetConfigHandling(key string, value any) {
	util.ErrorCheck(SetConfig(key, value), "Could not set config")
}

func SaveDeployedImage(tag, digest string) {
	deployedImage := model.PreviousImageOf(tag, digest)
	previous := GetPreviousDeployments()

	previous[arg_microserviceName] = append(previous[arg_microserviceName], deployedImage)
	util.Laugh(SetConfig("previous", previous))
}

func GetPreviousDeployments() PreviousDeployments {
	if config.Previous == nil {
		config.Previous = make(map[string]model.PreviousImages)
	}
	return config.Previous
}

func Registry() string {
	return config.Registry
}

func Repository() string {
	return config.Repository
}

func FullGcrRepositoryName() string {
	return config.Repository + GcrRepositoryName()
}

func GcrRepositoryName() string {
	mappings := config.Mappings[arg_microserviceName]
	if mappings.GCR != "" {
		util.Debug("Using GCR repository name from the configuration mappings: ", mappings.GCR)
		return mappings.GCR
	}
	return arg_microserviceName
}

func ResolveResourceName() string {
	return k8sNamespace + "-" + ContainerName()
}

// ContainerName returns container name from the command line argument.
// If the command line argument is set, it has a priority.
// Otherwise, it uses the mapping from the configuration file or the microservice name.
// Container name is also used as a part of the resource name.
func ContainerName() string {
	if arg_k8sResourceFullName != "" {
		util.Debug("Using container name from the command line flag '--k8s-name': ", arg_k8sResourceFullName)
		return arg_k8sResourceFullName
	}
	mappings := config.Mappings[arg_microserviceName]
	if mappings.K8S != "" {
		util.Debug("Using container name from the configuration file: ", mappings.K8S)
		return mappings.K8S
	}
	return arg_microserviceName
}

func ResolveResourceType() string {
	for _, statefulSet := range config.StatefulSets {
		if statefulSet == arg_microserviceName {
			return "statefulsets"
		}
	}
	return "deployments"
}
