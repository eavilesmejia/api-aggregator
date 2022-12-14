/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs server",
	Long: `Runs server.
	e.g: $ infura run server --restful
	`,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
