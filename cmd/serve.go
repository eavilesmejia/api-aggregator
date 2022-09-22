package cmd

import (
	"github.com/INFURA/infra/services/restful"
	"github.com/spf13/cobra"
)

var restOpt bool

// serverCmd represents the api command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Serves http restful service.",
	Long:  `Serves http restful service.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if restOpt {
			restful.Run()
		}

		return nil
	},
}

func init() {
	runCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVar(&restOpt, "restful", true, "Serves http restful service (default true)")
}
