package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	prompt "shumyk/kdeploy/cmd/prompt"
	util "shumyk/kdeploy/cmd/util"
)

func InitConfig(_ *cobra.Command, _ []string) {
	LoadConfiguration(nil, nil)
	validateVitalConfigs()
}

func LoadConfiguration(_ *cobra.Command, _ []string) {
	home, err := os.UserHomeDir()
	util.Laugh(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".kdeploy")
	viper.SetConfigType("yaml")

	_ = viper.SafeWriteConfig()
	util.Laugh(viper.ReadInConfig())
	util.Laugh(viper.Unmarshal(&config))
}

func validateVitalConfigs() {
	if len(config.Registry) == 0 {
		promptAndSaveConfig("registry")
	}
	if len(config.Repository) == 0 {
		promptAndSaveConfig("repository")
	}
}

func promptAndSaveConfig(configName string) {
	util.RedStderr(configName, " not found in ", viper.ConfigFileUsed())
	configValue, err := prompt.TextInput(configName)
	handleConfigPromptError(configName, err)
	SetConfigHandling(configName, configValue)
}

func handleConfigPromptError(configName string, err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			util.BoringStderr("Looks like you ctrl-c input. However, you can set it using:")
			util.BoringStderr(fmt.Sprintf("	kdeploy config set %v <value>", configName))
			util.BoringStderr("Or manually editing:")
			util.BoringStderr("	kdeploy config edit")
			os.Exit(1)
		}
		util.ErrorCheck(err, "Failed to request user input for missing configuration")
	}
}
