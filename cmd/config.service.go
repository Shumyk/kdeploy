package cmd

import (
	. "shumyk/kdeploy/cmd/model"
	. "shumyk/kdeploy/cmd/util"

	"github.com/spf13/viper"
)

func SetConfig(key string, value any) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func SetConfigHandling(key string, value any) {
	ErrorCheck(SetConfig(key, value), "Could not set config")
}

func SaveDeployedImage(tag, digest string) {
	deployedImage := PrevImageOf(tag, digest)
	previous := GetPreviousDeployments()

	previous[microservice] = append(previous[microservice], deployedImage)
	Laugh(SetConfig("previous", previous))
}

func GetPreviousDeployments() PreviousDeployments {
	if config.Previous == nil {
		config.Previous = make(map[string]PreviousImages)
	}
	return config.Previous
}

func Registry() string {
	return config.Registry
}

func Repository() string {
	return config.Repository
}

func BuildRepository(service string) string {
	return config.Repository + service
}

func ResolveResourceType(service string) string {
	for _, statefulSet := range config.StatefulSets {
		if statefulSet == service {
			return "statefulsets"
		}
	}
	return "deployments"
}
