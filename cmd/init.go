package cmd

import (
	"FloomCLI/config"
	"FloomCLI/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Define the initCmd
var initCmd = &cobra.Command{
	Use:   "init [deployment_type]",
	Short: "Initializes the Floom CLI by registering a new user if necessary",
	Long:  `This command checks for an existing API key and, if one isn't found, prompts for a deployment type before registering a new user with Floom.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var deploymentType string
		if len(args) > 0 {
			deploymentType = args[0]
		} else {
			// Prompt for deployment type if not provided as an argument
			fmt.Println("Please enter the deployment type (local, cloud, or custom endpoint):")
			fmt.Scanln(&deploymentType)
		}

		// Continue with the existing logic...
		initializeConfigForDeployment(deploymentType)
	},
}

func initializeConfigForDeployment(deploymentType string) {
	// Attempt to initialize configuration
	err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}

	// Check if API key already exists
	appConfig := config.GetConfig()
	if deployment, exists := appConfig.Deployments[deploymentType]; exists && deployment.Credentials.ApiKey != "" {
		fmt.Println("API key already exists for this deployment type. No need to register a new user.")
		return
	}

	// Validate deployment type
	if deploymentType != "local" && deploymentType != "cloud" {
		fmt.Println("Invalid deployment type. Please use 'local', 'cloud', or a valid custom endpoint.")
		return
	}

	// Register a new user
	registrationResponse, err := utils.RegisterUser(deploymentType)
	if err != nil {
		fmt.Printf("Error registering new user: %v\n", err)
		os.Exit(1)
	}

	// Update the configuration with new user details for the specified deployment type
	err = config.UpdateUserConfig(registrationResponse.ApiKey, registrationResponse.Username, registrationResponse.Nickname, deploymentType)
	if err != nil {
		fmt.Printf("Error updating configuration with new user details: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("New user registered and configuration updated successfully.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
