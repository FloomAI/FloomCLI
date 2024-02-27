package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type PipelineConfiguration struct {
	Name string `json:"name"`
	Url  string `json:"endpoint"`
	Port *int   `json:"port,omitempty"`
}

type DeploymentConfiguration struct {
	Credentials DeploymentCredentials   `json:"credentials"`
	Pipelines   []PipelineConfiguration `json:"pipelines"`
}

type DeploymentCredentials struct {
	ApiKey   string `json:"api_key"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// AppConfig holds the application configuration, including apiKey, username, and nickname.
type AppConfig struct {
	ConfigVersion string                             `json:"config_version"`
	Deployments   map[string]DeploymentConfiguration `json:"deployments"`
}

var (
	appConfig *AppConfig
	once      sync.Once
)

// GetConfig returns the instance of AppConfig.
func GetConfig() *AppConfig {
	if appConfig == nil {
		panic("config not initialized")
	}
	return appConfig
}

// GetConfigPath returns the path to the configuration directory for the app.
func GetConfigPath(appName string) (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	// Ensure the application-specific configuration directory exists.
	appConfigDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}

	return appConfigDir, nil
}

// InitConfig initializes the application configuration by loading it from a file.
// It uses GetConfigPath to determine the correct path for the configuration file.
// AppConfig and Pipeline are defined as before

// InitConfig initializes the application configuration by loading it from a file,
// or creates a new file with default configuration if it does not exist.
func InitConfig() error {
	var err error
	once.Do(func() {
		configPath, pathErr := GetConfigPath("floom-cli")
		if pathErr != nil {
			err = fmt.Errorf("failed to get config path: %v", pathErr)
			return
		}

		configFilePath := filepath.Join(configPath, "config.json")

		// Check if the config file exists, if not, create it with the default configuration.
		if _, openErr := os.Stat(configFilePath); os.IsNotExist(openErr) {
			defaultConfig := AppConfig{
				ConfigVersion: "1.0",
			}

			file, createErr := os.Create(configFilePath)
			if createErr != nil {
				err = fmt.Errorf("failed to create config file: %v", createErr)
				return
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			if encodeErr := encoder.Encode(defaultConfig); encodeErr != nil {
				err = fmt.Errorf("failed to encode default config: %v", encodeErr)
				return
			}

			// Since we just created the default config, set it as the current appConfig.
			appConfig = &defaultConfig
			return // Config is initialized with default, no need to load from file.
		}

		// Config file exists, proceed to open and decode it.
		file, openErr := os.Open(configFilePath)
		if openErr != nil {
			err = fmt.Errorf("failed to open config file: %v", openErr)
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		appConfig = &AppConfig{}
		if decodeErr := decoder.Decode(appConfig); decodeErr != nil {
			err = fmt.Errorf("failed to decode config file: %v", decodeErr)
			return
		}
	})

	return err
}

func (c *AppConfig) SaveConfig() error {
	configPath, err := GetConfigPath("floom-cli")
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}
	configFilePath := filepath.Join(configPath, "config.json")

	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create/open config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode and save config: %v", err)
	}

	return nil
}

func (c *AppConfig) AddOrUpdatePipeline(deploymentType, name, url string, port *int) {
	if c.Deployments == nil {
		c.Deployments = make(map[string]DeploymentConfiguration)
	}

	deploymentConfig, exists := c.Deployments[deploymentType]
	if !exists {
		deploymentConfig = DeploymentConfiguration{
			Credentials: DeploymentCredentials{}, // Initialize if necessary
			Pipelines:   []PipelineConfiguration{},
		}
	}

	// Check for an existing pipeline and update if found
	found := false
	for i, pipeline := range deploymentConfig.Pipelines {
		if pipeline.Name == name {
			found = true
			deploymentConfig.Pipelines[i] = PipelineConfiguration{Name: name, Url: url, Port: port}
			break
		}
	}

	// If not found, append a new pipeline
	if !found {
		newPipeline := PipelineConfiguration{Name: name, Url: url, Port: port}
		deploymentConfig.Pipelines = append(deploymentConfig.Pipelines, newPipeline)
	}

	// Update the deployment config in the app config
	c.Deployments[deploymentType] = deploymentConfig

	// Save the updated configuration
	if err := c.SaveConfig(); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
	}
}

func DeploymentConfigExists(deploymentType string) bool {
	_, exists := GetConfig().Deployments[deploymentType]
	return exists
}

// UpdateUserConfig function to update apiKey, username, and nickname in AppConfig.
func UpdateUserConfig(apiKey, username, nickname, deploymentType string) error {
	if appConfig == nil {
		return fmt.Errorf("config is not initialized")
	}

	// Check if the Deployments map is initialized; if not, initialize it.
	if appConfig.Deployments == nil {
		appConfig.Deployments = make(map[string]DeploymentConfiguration)
	}

	// Check if the specified deployment exists; if not, initialize it.
	if _, exists := appConfig.Deployments[deploymentType]; !exists {
		appConfig.Deployments[deploymentType] = DeploymentConfiguration{
			Credentials: DeploymentCredentials{},   // Initialize with empty credentials
			Pipelines:   []PipelineConfiguration{}, // Initialize with an empty slice
		}
	}

	// Update the credentials for the specific deployment type.
	deploymentConfig := appConfig.Deployments[deploymentType]
	deploymentConfig.Credentials = DeploymentCredentials{
		ApiKey:   apiKey,
		Username: username,
		Nickname: nickname,
	}

	// Important: Update the map entry with the modified deploymentConfig
	appConfig.Deployments[deploymentType] = deploymentConfig

	// Use the existing SaveConfig function to save the updated AppConfig.
	return appConfig.SaveConfig()
}

// GetApiKeyForDeployment returns the API key for a given deployment type.
func GetApiKeyForDeployment(deploymentType string) (string, error) {
	appConfig := GetConfig() // Assuming GetConfig() fetches the current AppConfig instance.

	// Check if the deployment exists in the appConfig
	deploymentConfig, exists := appConfig.Deployments[deploymentType]
	if !exists {
		return "", fmt.Errorf("deployment type '%s' not found in configuration", deploymentType)
	}

	apiKey := deploymentConfig.Credentials.ApiKey
	if apiKey == "" {
		return "", fmt.Errorf("API key for deployment type '%s' is missing", deploymentType)
	}

	return apiKey, nil
}
