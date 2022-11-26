package cmd

import (
	"fmt"
	"shumyk/kdeploy/cmd/model"
	util "shumyk/kdeploy/cmd/util"

	"github.com/AlecAivazis/survey/v2"
)

const (
	imageSelectTitle = "select image to deploy"
	repoSelectTitle  = "select repo"
)

func TextInput(name string) (result string, err error) {
	textInput := survey.Input{
		Message: fmt.Sprintf("please enter %v", name),
	}
	err = survey.AskOne(&textInput, &result)
	return
}

func ImageSelect(input model.PromptInputs) model.SelectedImage {
	options := input.ImageOptions().Sorted().Stringify()
	chosenString := prompt(imageSelectTitle, options)
	return model.ParseSelectedImage(chosenString)
}

func RepoSelect(repos []string) string {
	return prompt(repoSelectTitle, repos)
}

func prompt(title string, options []string) (selected string) {
	prompt := &survey.Select{
		Message: title,
		Options: options,
	}
	err := survey.AskOne(prompt, &selected)

	util.Laugh(err)
	util.TerminateOnSigint(selected)
	return
}
