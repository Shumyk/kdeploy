package cmd

import (
	. "shumyk/kdeploy/cmd/model"
	. "shumyk/kdeploy/cmd/util"
)

var config configuration

type configuration struct {
	Registry     string              `yaml:"registry,omitempty"`
	Repository   string              `yaml:"repository,omitempty"`
	StatefulSets []string            `yaml:"statefulSets,omitempty"`
	Previous     PreviousDeployments `yaml:"previous,omitempty" conf:"no"`
}

type PreviousDeployments map[string]PreviousImages

func (c configuration) View() *configuration {
	c.Previous = nil
	return &c
}

func (p PreviousDeployments) Keys() []string {
	keyMapping := ReturnKey[string, PreviousImages]
	return MapToSliceMapping(p, keyMapping)
}
