package cmd

import (
	prompt "shumyk/kdeploy/cmd/prompt"
	. "shumyk/kdeploy/cmd/util"
)

func KDeploy() {
	DeployNew()
}

func KDeployWithRegistry() {
	repos := ListRepos()
	microservice = prompt.RepoSelect(repos)
	DeployNew()
}

func KDeployPrevious() {
	previous := GetPreviousDeployments()[microservice]
	TerminateOnEmpty(previous, "previous deployments of", microservice, "absent")
	DeployPrevious(previous)
}

func KDeployPreviousWithRegistry() {
	repos := GetPreviousDeployments().Keys()
	TerminateOnEmpty(repos, "previous deployments absent")

	microservice = prompt.RepoSelect(repos)
	KDeployPrevious()
}
