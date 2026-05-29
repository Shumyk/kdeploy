package cmd

import (
	util "shumyk/kdeploy/cmd/util"

	"github.com/spf13/cobra"
)

func DestroyContext(_ *cobra.Command, _ []string) error {
	if garGuy == nil || garGuy.client == nil {
		return nil
	}

	util.Debug("Destroying GAR client")
	return garGuy.client.Close()
}
