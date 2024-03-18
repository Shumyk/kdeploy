package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
	configDir := home + "/.config/kdeploy"
	configName := ".kdeploy"
	configType := "yaml"
	configPath := filepath.Join(configDir, configName+"."+configType)

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	createConfigFileIfNotExists(configDir, configPath)
	util.Laugh(viper.ReadInConfig())
	util.Laugh(viper.Unmarshal(&config))
}

func createConfigFileIfNotExists(configDir, configPath string) {
	if _, err := os.Stat(configPath); os.IsExist(err) {
		return
	}
	err := os.MkdirAll(configDir, os.ModePerm)
	util.ErrorCheck(err, "Failed to create config directory")
	file, err := os.Create(configPath)
	util.ErrorCheck(err, "Failed to create config file")
	file.Close()
}

func validateVitalConfigs() {
	if len(config.Registry) == 0 {
		promptAndSaveConfig("registry", "*gcr.io")
	}
	if len(config.Repository) == 0 {
		promptAndSaveConfig("repository", "your-domain-infra/domain/domain-")
	}
}

func promptAndSaveConfig(configName, example string) {
	util.PurpleStout(configName, " not found in ", viper.ConfigFileUsed())
	configValue, err := prompt.TextInput(configName, example)
	handleConfigPromptError(configName, err)
	SetConfigHandling(configName, configValue)
}

func handleConfigPromptError(configName string, err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			util.PurpleStout("did you ctrl-c me? anyway, you can set it using:")
			util.BoringStderr(fmt.Sprintf("\tkdeploy config set %v <value>", configName))
			util.BoringStderr("or manually editing:")
			util.BoringStderr("\tkdeploy config edit")
			os.Exit(1)
		}
		util.ErrorCheck(err, "Failed to request user input for missing configuration")
	}
}
