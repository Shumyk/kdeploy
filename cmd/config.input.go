package cmd

import (
	"os"
	prompt "shumyk/kdeploy/cmd/prompt"
	util "shumyk/kdeploy/cmd/util"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/viper"
)

func inputVitalConfig(configName, example string) {
	util.PurpleStout(configName, " not found in ", viper.ConfigFileUsed())
	configValue := inputConfig(configName, example, true)
	SetConfigHandling(configName, configValue)
}

func inputConfig(configName, example string, retry bool) (value string) {
	for {
		var err error
		value, err = prompt.TextInput(configName, example)
		handleConfigPromptError(err)

		if retry && len(value) == 0 {
			util.CyanStout("looks like empty input, try to press some buttons this time")
		} else {
			return
		}
	}
}

func handleConfigPromptError(err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			util.DashLine()
			util.PurpleStout("did you ctrl-c me? anyway, you can set config using:")
			util.BoringStderr("\tkdeploy config set <property> <value>")
			util.BoringStderr("or define complex configs interactively:")
			util.BoringStderr("\tkdeploy config define mappings")
			util.BoringStderr("or manually editing:")
			util.BoringStderr("\tkdeploy config edit")
			os.Exit(1)
		}
		util.ErrorCheck(err, "Failed to request user input for missing configuration")
	}
}
