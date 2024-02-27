package cmd

import (
	"FloomCLI/config"
	"FloomCLI/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// deployCmd represents the deployment command
var deployCmd = &cobra.Command{
	Use:   "deploy [local|cloud|custom_endpoint] [file]",
	Short: "Deploy pipeline configurations to Floom",
	Long: `Deploy pipeline configurations to a local Floom Docker instance or to the Floom cloud.

For local deployment, use:
    floom deploy local pipeline.yml

For cloud deployment, use:
    floom deploy cloud pipeline.yml

For custom endpoint deployment, use:
	floom deploy http://184.152.3.12 pipeline.yml`,

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentType := args[0]
		yamlFile := args[1]

		// validate deploymentType is non empty and yamlFile are valid

		if deploymentType == "" {
			fmt.Println("Deployment type is required, use 'local' or 'cloud'.")
			return
		}

		if yamlFile == "" {
			fmt.Println("YAML file is required")
			return
		}

		deploy(deploymentType, yamlFile)
	},
}

func resolveYamlPath(yamlFile string, cwd ...string) (string, error) {
	if filepath.IsAbs(yamlFile) {
		return yamlFile, nil
	}

	workingDir := ""
	if len(cwd) > 0 {
		workingDir = cwd[0]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		workingDir = dir
	}

	// Normalize paths for consistent comparisons
	normalizedYamlFile := filepath.Clean(yamlFile)
	normalizedWD := filepath.Clean(workingDir)

	// Check if it refers to the current directory (case-insensitively)
	if strings.EqualFold(normalizedYamlFile, normalizedWD) {
		return yamlFile, nil
	}

	// Prepend the current directory explicitly
	return filepath.Join(workingDir, yamlFile), nil
}

func deploy(deploymentType string, yamlFile string) {
	// Implementation for deploying to local Floom Docker instance
	appConfig := config.GetConfig()

	// Check for cloud deployment configuration; initialize if not found
	if deploymentType == "cloud" && !config.DeploymentConfigExists(deploymentType) {
		fmt.Println("Cloud deployment configuration not found. Initializing...")
		initializeConfigForDeployment(deploymentType)
	}

	yamlPath, err := resolveYamlPath(yamlFile)
	if err != nil {
		fmt.Println("Error resolving YAML path:", err)
		return
	}

	// 1. Parse YAML to find context files
	FloomYaml, err := utils.ParseYaml(yamlPath)
	if err != nil {
		fmt.Println("Error parsing Floom YAML file:", err)
		return
	}

	// 2. Upload context files and get asset IDs

	// implement FloomYaml.Prompt.Context iteration
	for _, context := range FloomYaml.Pipeline.Prompt.Context {
		pathInterface := context.Configuration["path"]

		// Check if pathInterface is a string (single path)
		if path, ok := pathInterface.(string); ok {
			fileId, err := utils.UploadFile(deploymentType, path)
			if err != nil {
				fmt.Println("Error uploading the file:", err)
				return
			}
			// Replace 'path' with 'assetId' (array of one element)
			context.Configuration["assetId"] = []string{fileId}
		} else if paths, ok := pathInterface.([]interface{}); ok {
			// Handle case where path is an array of strings
			var fileIds []string
			for _, pathElement := range paths {
				path, ok := pathElement.(string)
				if !ok {
					fmt.Println("Error: path is not a string in the array")
					return
				}
				fileId, err := utils.UploadFile(deploymentType, path)
				if err != nil {
					fmt.Println("Error uploading the file:", err)
					return
				}
				fileIds = append(fileIds, fileId)
			}
			// Replace 'path' with 'assetId' (array of file IDs)
			context.Configuration["assetId"] = fileIds
		}
		// Remove the original 'path' entry
		delete(context.Configuration, "path")
	}

	// 3. Replace context paths with asset IDs in the YAML
	// 4. Commit the modified pipeline configuration
	err = utils.DeployPipeline(deploymentType, *FloomYaml)
	if err != nil {
		fmt.Println("Error deploying pipeline:", err)
		return
	}

	// Fetch the API key and username for the deployment
	apiKey, err := config.GetApiKeyForDeployment(deploymentType)
	if err != nil {
		fmt.Println("Error fetching API key for deployment:", err)
		return
	}

	// Assuming GetConfig() and Username retrieval based on updated config structure
	deploymentConfig, exists := appConfig.Deployments[deploymentType]
	if !exists {
		fmt.Println("Deployment type not found in configuration")
		return
	}
	username := deploymentConfig.Credentials.Username

	// Construct the pipeline URL
	pipelineURL := fmt.Sprintf("https://%s-%s.pipeline.floom.ai/", FloomYaml.Pipeline.Name, username)

	// Print success message, pipeline URL, and instructions for making HTTP POST request
	fmt.Printf("Pipeline '%s' deployed successfully.\n", FloomYaml.Pipeline.Name)
	fmt.Println("Pipeline URL:", pipelineURL)
	fmt.Println("You can send an HTTP POST request to this URL with the following headers:")
	fmt.Println("API-Key:", apiKey)
	fmt.Println("Content-Type: application/json")
	fmt.Println("In the request body, include a JSON with a 'prompt' field, for example:")
	fmt.Println(`{"prompt": "Your prompt example here"}`)

	// Construct the pipeline URL for cloud deployments
	var port *int // Set to nil by default, indicating cloud deployment or irrelevant port
	if deploymentType != "cloud" {
		// For non-cloud deployments, you might have a specific port to use
		// Example: port = new(int); *port = 8080
	}

	// Add the pipeline to the configuration
	appConfig.AddOrUpdatePipeline(deploymentType, FloomYaml.Pipeline.Name, pipelineURL, port)

}

func init() {
	rootCmd.AddCommand(deployCmd)
	// Here you can define flags and configuration settings.
}
