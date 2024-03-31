package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	createConfigFileIfNotExists(configDir, configPath)
	util.Laugh(viper.ReadInConfig(), "Failed to read configuration file")
	util.Laugh(viper.Unmarshal(&config), "Failed to unmarshal configuration")

	initContext()
}

func createConfigFileIfNotExists(configDir, configPath string) {
	if _, err := os.Stat(configPath); nil == err || os.IsExist(err) {
		util.Debug("Config file exists")
		return
	}
	util.Debug("Creating config file")
	err := os.MkdirAll(configDir, os.ModePerm)
	util.ErrorCheck(err, "Failed to create config directory")
	file, err := os.Create(configPath)
	util.ErrorCheck(err, "Failed to create config file")
	file.Close()
}

func validateVitalConfigs() {
	if len(config.Registry) == 0 {
		inputVitalConfig("registry", "*gcr.io")
	}
	if len(config.Repository) == 0 {
		inputVitalConfig("repository", "your-domain-infra/domain/domain-")
	}
}

func initContext() {
	util.SetDebugMode(config.Debug)
	util.Debug("Debug mode is enabled")
}
