package cmd

import (
	"FloomCLI/utils"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [path to docker-compose file]",
	Short: "Stops the Floom environment using Docker Compose",
	Long: `Stops the Floom environment by automatically detecting and using a Docker Compose file in the current directory,
or by using a Docker Compose file specified as an argument.

This command searches for 'docker-compose.yml' in the current directory to use as the default Docker Compose file. If no such file is found, or if a specific path is provided as an argument, that file is used instead.

Examples:
  # Stop using the default docker-compose.yml in the current directory
  floom stop

  # Stop using a specific Docker Compose file
  floom stop path/to/your/docker-compose.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultComposeFile := "docker-compose.yml"
		composeFilePath := ""

		// Check if a file path is provided as an argument
		if len(args) > 0 {
			composeFilePath = args[0]
		} else {
			// Search for the default Docker Compose file in the current directory
			files, err := filepath.Glob("./*.yml")
			if err != nil {
				log.Fatalf("Error reading directory: %v", err)
			}
			for _, file := range files {
				if filepath.Base(file) == defaultComposeFile {
					composeFilePath = file
					break
				}
			}
		}

		// Check if a Docker Compose file has been identified or provided
		if composeFilePath == "" {
			log.Fatalf("No Docker Compose file found or provided")
		} else {
			log.Printf("Using Docker Compose file: %s", composeFilePath)
		}

		// Stop the Docker Compose services
		log.Println("Stopping Docker Compose services...")
		_, err := utils.ExecuteShellCommand("docker-compose", []string{"-f", composeFilePath, "down"})
		if err != nil {
			log.Fatalf("Error stopping Docker Compose services: %v", err)
		}
		log.Println("Docker Compose stopped successfully.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
