package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	prompt "shumyk/kdeploy/cmd/prompt"
	. "shumyk/kdeploy/cmd/util"
)

func InitConfig(_ *cobra.Command, _ []string) {
	LoadConfiguration(nil, nil)
	validateVitalConfigs()
}

func LoadConfiguration(_ *cobra.Command, _ []string) {
	home, err := os.UserHomeDir()
	Laugh(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".kdeploy")
	viper.SetConfigType("yaml")

	_ = viper.SafeWriteConfig()
	Laugh(viper.ReadInConfig())
	Laugh(viper.Unmarshal(&config))
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
	RedStderr(configName, " not found in ", viper.ConfigFileUsed())
	configValue, err := prompt.TextInput(configName)
	handleConfigPromptError(configName, err)
	SetConfigHandling(configName, configValue)
}

func handleConfigPromptError(configName string, err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			BoringStderr("Looks like you ctrl-c input. However, you can set it using:")
			BoringStderr(fmt.Sprintf("	kdeploy config set %v <value>", configName))
			BoringStderr("Or manually editing:")
			BoringStderr("	kdeploy config edit")
			os.Exit(1)
		}
		ErrorCheck(err, "Failed to request user input for missing configuration")
	}
}
