package cmd

import (
	prompt "shumyk/kdeploy/cmd/prompt"
	util "shumyk/kdeploy/cmd/util"
)

func KDeploy() {
	DeployNew()
}

func KDeployWithRegistry() {
	repos := ListRepos()
	arg_microserviceName = prompt.RepoSelect(repos)
	DeployNew()
}

func KDeployPrevious() {
	previous := GetPreviousDeployments()[arg_microserviceName]
	util.TerminateOnEmpty(previous, "previous deployments of", arg_microserviceName, "absent")
	DeployPrevious(previous)
}

func KDeployPreviousWithRegistry() {
	var repos []string = GetPreviousDeployments().Keys()
	util.TerminateOnEmpty(repos, "previous deployments absent")

	arg_microserviceName = prompt.RepoSelect(repos)
	KDeployPrevious()
}
