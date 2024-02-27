package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

const floomAIHelpArt = `
  ______ _                                  _____ 
 |  ____| |                           /\   |_   _|
 | |__  | | ___   ___  _ __ ___      /  \    | |  
 |  __| | |/ _ \ / _ \| '_   _ \    / /\ \   | |
 | |    | | (_) | (_) | | | | | |  / ____ \ _| |_
 |_|    |_|\___/ \___/|_| |_| |_| /_/    \_\_____|`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "floom",
	Short: "FloomCLI is a command-line interface for managing Floom environments.",
	Long: `FloomCLI provides a comprehensive suite of tools for configuring, managing, and interacting
with Floom environments directly from the command line. 
...
(command-line toolset for the efficient management of Floom environments.)`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, strings []string) {
		color.New(color.FgHiCyan).Fprintln(os.Stdout, floomAIHelpArt)
		fmt.Println()
		cmd.Usage()
	})

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}
