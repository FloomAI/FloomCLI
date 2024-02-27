package cmd

import (
	"FloomCLI/config" // Import your config package
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays the configuration file path and its content",
	Long:  `This command fetches and displays the configuration file path along with its content.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayConfigInfo()
	},
}

func displayConfigInfo() {
	// Assuming GetConfigPath is a function that returns the full path of the config file
	configDirPath, err := config.GetConfigPath("floom-cli") // Adjust the argument as per your actual config path function
	if err != nil {
		fmt.Printf("Error fetching config path: %v\n", err)
		return
	}
	configFilePath := filepath.Join(configDirPath, "config.json")
	fmt.Println("Configuration File Path:", configFilePath)

	// Read and display the content of the config file using os.ReadFile
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}
	fmt.Println("Configuration Content:")
	fmt.Println(string(content))
}

func init() {
	rootCmd.AddCommand(infoCmd) // Make sure your rootCmd is correctly initialized as per Cobra setup
}
