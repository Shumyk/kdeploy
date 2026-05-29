package cmd

import (
	"maps"
	util "shumyk/kdeploy/cmd/util"
	"slices"

	"github.com/spf13/cobra"
)

var (
	arg_microserviceName    string
	arg_previousMode        bool
	arg_k8sResourceFullName string

	kdeploy = cobra.Command{
		Use:   "kdeploy [microservice]",
		Short: "k[8s]deploy - deploy from the terminal",
		Long: `Searches for images of requested microservice in Google Artifact Registry,
Prompts you to interactively select an image for deployment (arrows navigation, search features),
And sets the selected image in the workload.
If microservice was not specified - it obtains possible repositories from the registry and prompts you to select it first.

kdeploy requires GAR and K8S configuration blocks.
Define them interactively:
    kdeploy config define gar
    kdeploy config define k8s
Or edit them manually:
    kdeploy config edit

Assumed that all workloads are of Deployment type. If some are StatefulSets, set them in configurations:
    kdeploy config set statefulsets ms-events,ms-core

kdeploy remembers every deployment you made and allows you to redeploy previous images.
    kdeploy --previous [microservice]`,
		Args:               cobra.MaximumNArgs(1),
		PreRun:             InitConfig,
		PersistentPostRunE: DestroyContext,
		Run:                kdeployRun,
	}

	// configurations commands
	configCmd = cobra.Command{
		Use:              "config [action] [args]...",
		Short:            "View, edit, set configurations",
		PersistentPreRun: LoadConfiguration,
	}
	configViewCmd = cobra.Command{
		Use:   "view",
		Short: "Displays current configuration.",
		Run:   runConfigView,
		Args:  cobra.NoArgs,
	}
	configEditCmd = cobra.Command{
		Use:   "edit",
		Short: "Basically, opens vim editor on configuration file.",
		Run:   RunConfigEdit,
		Args:  cobra.NoArgs,
	}
	configSetCmd = cobra.Command{
		Use:   "set [property] [value]",
		Short: "Conveniently set properties.",
		Long: `Conveniently set properties.
Use ',' delimiter (without space) for array type properties (e.x. statefulsets).`,
		Example: `  kdeploy config set debug true
  kdeploy config set statefulsets ms-events,ms-core`,
		Run:  RunConfigSet,
		Args: cobra.ExactArgs(2),
	}
	configDefineCmd = cobra.Command{
		Use:   "define [property]",
		Short: "Define complex property in configuration file.",
		Long: `Define complex properties in configuration file.
You will be prompted to enter values interactively.

Currently supported complex properties:
  gar
    Defines Google Artifact Registry settings:
    project, location, repository and optional package prefix.
    ref: projects/{project}/locations/{location}/repositories/{repo}/packages/{package-prefix}{cli-argument}

  k8s
    Defines the image path used in Kubernetes manifests:
    registry host and repository prefix.
    ref: {registry}/{repository}{gar-package}:{tag}@sha256:{digest}

  mappings
    Defines custom service mappings for a microservice:
    service name, GAR package name and K8S resource/container name.
    Use this when GAR or K8S names differ from the microservice name,
    so you do not have to pass --k8s-name every time.`,
		Example: `  kdeploy config define gar
  kdeploy config define k8s
  kdeploy config define mappings`,
		Run:       RunConfigDefine,
		Args:      cobra.ExactArgs(1),
		ValidArgs: slices.Collect(maps.Keys(complexConfigurations)),
	}
)

func kdeployRun(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		deploySelectingRegistry()
	} else {
		deployMicroservice(args)
	}
}

func deploySelectingRegistry() {
	if arg_previousMode {
		KDeployPreviousWithRegistry()
	} else {
		KDeployWithRegistry()
	}
}

func deployMicroservice(args []string) {
	arg_microserviceName = args[0]
	util.Debug("Deploying microservice: ", arg_microserviceName)
	if arg_previousMode {
		KDeployPrevious()
	} else {
		KDeploy()
	}
}

func Execute() {
	util.ErrorCheck(kdeploy.Execute(), "Failed to execute kdeploy :|")
}

func init() {
	kdeploy.Flags().BoolVarP(&arg_previousMode, "previous", "p", false, "deploy previous images")
	kdeploy.Flags().StringVarP(&arg_k8sResourceFullName, "k8s-name", "n", "", "k8s name to use for deployment")

	configCmd.AddCommand(&configViewCmd, &configSetCmd, &configDefineCmd, &configEditCmd)
	kdeploy.AddCommand(&configCmd)
}
