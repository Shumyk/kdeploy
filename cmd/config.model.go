package cmd

import (
	model "shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"
)

var config configuration

type configuration struct {
	Debug bool `yaml:"debug,omitempty"`

	GAR          GAR                        `yaml:"gar,omitempty"`
	K8S          K8S                        `yaml:"k8s,omitempty"`
	StatefulSets []string                   `yaml:"statefulSets,omitempty"`
	Mappings     map[string]ServiceMappings `yaml:"mappings,omitempty"`
	Previous     PreviousDeployments        `yaml:"previous,omitempty" conf:"no"`
}

type GAR struct {
	Project       string `yaml:"project"`
	Location      string `yaml:"location"`
	Repository    string `yaml:"repository"`
	PackagePrefix string `yaml:"packagePrefix,omitempty"`
}

func (g GAR) isValid() bool {
	if g.Project == "" {
		return false
	}
	if g.Location == "" {
		return false
	}
	if g.Repository == "" {
		return false
	}

	return true
}

type K8S struct {
	Registry   string `yaml:"registry"`
	Repository string `yaml:"repository"`
}

func (k K8S) isValid() bool {
	if k.Registry == "" {
		return false
	}
	if k.Repository == "" {
		return false
	}

	return true
}

type ServiceMappings struct {
	GAR string `yaml:"gar,omitempty"`
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
