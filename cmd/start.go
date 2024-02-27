package cmd

import (
	"FloomCLI/utils"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [path to docker-compose file]",
	Short: "Starts the Floom environment using Docker Compose",
	Long: `Starts the Floom environment by automatically detecting and using a Docker Compose file in the current directory,
or by using a Docker Compose file specified as an argument.

This command searches for 'docker-compose-local.yml' in the current directory to use as the default Docker Compose file. If no such file is found, or if a specific path is provided as an argument, that file is used instead.

Upon successfully launching Docker Compose, the command also performs additional initialization steps for Floom, such as configuring settings and finding the local IP address of the Floom Docker container.

Examples:
  # Start using the default docker-compose.yml in the current directory
  floom start

  # Start using a specific Docker Compose file
  floom start path/to/your/docker-compose.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		// Define the default Docker Compose file name
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

		// Check if Docker Compose services are running and stop them if they are
		log.Println("Checking if Docker Compose services are already running...")
		checkOutput, checkErr := utils.ExecuteShellCommand("docker-compose", []string{"-f", composeFilePath, "ps", "-q"})
		if checkErr != nil {
			log.Fatalf("Error checking Docker Compose services: %v", checkErr)
		}
		if checkOutput != "" {
			log.Println("Stopping running Docker Compose services...")
			_, downErr := utils.ExecuteShellCommand("docker-compose", []string{"-f", composeFilePath, "down"})
			if downErr != nil {
				log.Fatalf("Error stopping Docker Compose services: %v", downErr)
			}
		}

		// Restart the Docker Compose services
		log.Printf("Starting Docker Compose with file: %s", composeFilePath)
		output, err := utils.ExecuteShellCommand("docker-compose", []string{"-f", composeFilePath, "-p", "floom", "up", "-d"})

		if err != nil {
			log.Fatalf("Error executing Docker Compose: %v\n%s", err, output)
		}
		log.Println("Docker Compose started successfully.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
