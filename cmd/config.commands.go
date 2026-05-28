package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	util "shumyk/kdeploy/cmd/util"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	COMPLEX_MAPPINGS = "mappings"
	COMPLEX_GAR      = "gar"
	COMPLEX_K8S      = "k8s"

	complexConfigurations = map[string]struct{}{
		COMPLEX_MAPPINGS: {},
		COMPLEX_GAR:      {},
		COMPLEX_K8S:      {},
	}
)

func runConfigView(_ *cobra.Command, _ []string) {
	viewBytes, err := yaml.Marshal(config.View())
	util.ErrorCheck(err, "Failed to marshal configuration to yaml")
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
	util.RedStderr("Non existing property: ", property)
	util.BoringStderr("Possible configuration properties:")
	util.ErrorCheck(properties.Flush(), "Could not print configuration properties")
}

func RunConfigDefine(_ *cobra.Command, args []string) {
	property := args[0]
	if _, ok := complexConfigurations[property]; !ok {
		util.RedStderr("Non existing complex property: ", property)
		return
	}

	switch property {
	case COMPLEX_MAPPINGS:
		handleDefineMappings()
	case COMPLEX_GAR:
		handleDefineGar()
	case COMPLEX_K8S:
		handleDefineK8s()
	}
}

func handleDefineMappings() {
	serviceName := inputConfig("service name", "api-events", true)
	gcr := inputConfig("GAR", "events", false)
	k8s := inputConfig("K8S", "cmpn-events", false)

	config.Mappings[serviceName] = ServiceMappings{GAR: gcr, K8S: k8s}
	SetConfigHandling("mappings", config.Mappings)
}

func handleDefineGar() {
	config.GAR = GAR{
		Project:       inputConfig("gcloud project", "company-infra", true),
		Location:      inputConfig("location", "us", true),
		Repository:    inputConfig("repository", "docker-images", true),
		PackagePrefix: inputConfig("packages prefix, optional", "company-", false),
	}
	SetConfigHandling("gar", config.GAR)
}

func handleDefineK8s() {
	config.K8S = K8S{
		Registry:   inputConfig("registry", "us-docker.pkg.dev", true),
		Repository: inputConfig("repository", "google-infra/docker/", true),
	}
	SetConfigHandling("k8s", config.K8S)
}

func RunConfigEdit(_ *cobra.Command, _ []string) {
	vim := exec.Command("vim", viper.ConfigFileUsed())
	vim.Stdin, vim.Stdout = os.Stdin, os.Stdout
	util.ErrorCheck(vim.Run(), "Error editing configuration")
}
