package cmd

import (
	"github.com/spf13/cobra"
	. "shumyk/kdeploy/cmd/util"
)

var (
	previousMode bool

	kdeploy = cobra.Command{
		Use:   "kdeploy [microservice]",
		Short: "k[8s]deploy - deploy from the terminal",
		Long: `Searches for images of requested microservice in Google Container Registry,
Prompts you to interactively select an image for deployment (arrows navigation, search features),
And sets the selected image in the workload.
If microservice was not specified - it obtains possible repositories from the registry and prompts you to select it first.

kdeploy requires two configuration properties - registry and repository.
The registry is where to look for your images (e.x. us.gcr.io), and the repository is the path to your images.
Set them using:
    kdeploy config set [registry|repository] [value]
Or  kdeploy config edit

Assumed that all workloads are of Deployment type. If some are StatefulSets, set them in configurations:
    kdeploy config set statefulsets ms-events,ms-core

kdeploy remembers every deployment you made and allows you to redeploy previous images.
    kdeploy --previous [microservice]`,
		Args:   cobra.MaximumNArgs(1),
		PreRun: InitConfig,
		Run:    kdeployRun,
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
		Example: `  kdeploy config set registry us.gcr.io  
  kdeploy config set statefulsets ms-events,ms-core`,
		Run:  RunConfigSet,
		Args: cobra.ExactArgs(2),
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
	if previousMode {
		KDeployPreviousWithRegistry()
	} else {
		KDeployWithRegistry()
	}
}

func deployMicroservice(args []string) {
	microservice = args[0]
	if previousMode {
		KDeployPrevious()
	} else {
		KDeploy()
	}
}

func Execute() {
	ErrorCheck(kdeploy.Execute(), "Failed to execute kdeploy :|")
}

func init() {
	kdeploy.Flags().BoolVarP(&previousMode, "previous", "p", false, "deploy previous images")

	configCmd.AddCommand(&configViewCmd, &configSetCmd, &configEditCmd)
	kdeploy.AddCommand(&configCmd)
}
