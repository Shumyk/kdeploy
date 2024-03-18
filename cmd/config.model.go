package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"
)

var config configuration

type configuration struct {
	Registry     string              `yaml:"registry,omitempty"`
	Repository   string              `yaml:"repository,omitempty"`
	StatefulSets []string            `yaml:"statefulSets,omitempty"`
	Mappings     map[string]Names    `yaml:"mappings,omitempty"`
	Previous     PreviousDeployments `yaml:"previous,omitempty" conf:"no"`
}

type Names struct {
	GCR string `yaml:"gcr,omitempty"`
	K8S string `yaml:"k8s,omitempty"`
}

type PreviousDeployments map[string]model.PreviousImages

func (c configuration) View() *configuration {
	c.Previous = nil
	return &c
}

func (p PreviousDeployments) Keys() []string {
	keyMapping := util.ReturnKey[string, model.PreviousImages]
	return util.MapToSliceMapping(p, keyMapping)
}
