package cmd

import (
	"FloomCLI/config"
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the Floom CLI",
	Long:  `This command displays the current version of the Floom CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		DisplayVersion() // Call the function to display the version
	},
}

func DisplayVersion() {
	err := config.InitConfig() // This will initialize the config if it doesn't exist.
	if err != nil {
		fmt.Printf("Error initializing or loading the configuration: %v\n", err)
		return
	}

	config := config.GetConfig() // Now we are sure that config is initialized.
	fmt.Printf("Floom CLI Version: %s\n", config.ConfigVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd) // Add the versionCmd to the root command
}
