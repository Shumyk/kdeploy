package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"reflect"
	. "shumyk/kdeploy/cmd/util"
	"strings"
	"text/tabwriter"
)

func runConfigView(_ *cobra.Command, _ []string) {
	viewBytes, err := yaml.Marshal(config.View())
	ErrorCheck(err, "Couldn't marshal config file")
	fmt.Println(string(viewBytes))
}

func RunConfigSet(_ *cobra.Command, args []string) {
	property, valueRaw := args[0], args[1]
	properties := tabwriter.NewWriter(os.Stderr, 1, 2, 4, ' ', tabwriter.TabIndent)

	configValue := reflect.ValueOf(&config)
	for i := 0; i < configValue.Elem().NumField(); i++ {
		field := configValue.Elem().Type().Field(i)
		// FIXME: this won't work if we have more than 1 value for tag
		if field.Tag.Get("conf") == "no" {
			continue
		}
		if strings.EqualFold(field.Name, property) {
			var value any = valueRaw
			if field.Type.Kind() == reflect.Slice {
				value = strings.Split(valueRaw, ",")
			}
			SetConfigHandling(field.Name, value)
			return
		}
		_, _ = fmt.Fprintln(properties, "\t"+strings.ToLower(field.Name)+"\t:\t"+field.Type.String())
	}
	RedStderr("Non existing property: ", property)
	BoringStderr("Possible configuration properties:")
	ErrorCheck(properties.Flush(), "Could not print configuration properties")
}

func RunConfigEdit(_ *cobra.Command, _ []string) {
	vim := exec.Command("vim", viper.ConfigFileUsed())
	vim.Stdin, vim.Stdout = os.Stdin, os.Stdout
	ErrorCheck(vim.Run(), "Error editing configuration")
}
