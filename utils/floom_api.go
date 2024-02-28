package utils

import (
	"FloomCLI/config"
	"FloomCLI/models"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileUploadResponse struct {
	FileId string `json:"fileId"`
}

type UserRegistrationResponse struct {
	ApiKey   string `json:"apiKey"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func getBaseUrl(deploymentType string) string {
	if deploymentType == "local" || deploymentType == "localhost" {
		return "http://localhost:4050"
	}

	if deploymentType == "cloud" {
		return "https://api.floom.ai"
	}

	return deploymentType
}

// RegisterUser sends a request to register a new user and returns the API key, username, and nickname.
func RegisterUser(deploymentType string) (UserRegistrationResponse, error) {
	var registrationResponse UserRegistrationResponse

	url := getBaseUrl(deploymentType) + "/v1/Users/Register"
	// The example doesn't specify body content, assuming no payload is required for registration.
	// If a payload is needed, construct it here.
	requestBody, err := json.Marshal(map[string]string{
		// Add required fields for the request body if needed.
	})
	if err != nil {
		return registrationResponse, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return registrationResponse, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return registrationResponse, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return registrationResponse, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&registrationResponse); err != nil {
		return registrationResponse, fmt.Errorf("error decoding response: %v", err)
	}

	return registrationResponse, nil
}

func UploadFile(deploymentType, filePath string) (string, error) {

	var apiKey string
	var err error

	if deploymentType == "cloud" {
		apiKey, err = config.GetApiKeyForDeployment(deploymentType)
		if err != nil {
			return "", err
		}
	}

	// Open the file to be sent
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Create the form file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	// Copy the file data to the form file
	if _, err := io.Copy(fileWriter, file); err != nil {
		return "", fmt.Errorf("error copying file data: %w", err)
	}

	// Close the multipart writer to finalize the multipart body
	if err := multiPartWriter.Close(); err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	url := getBaseUrl(deploymentType) + "/v1/Assets"
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set the content type header, including the boundary
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	// Set the API key header
	if apiKey != "" {
		req.Header.Set("Api-Key", apiKey)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	// Decode the JSON response
	var uploadResponse FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
	return uploadResponse.FileId, nil
}

// DeployPipeline Deploy sends the pipeline to the Floom API for deployment
func DeployPipeline(deploymentType string, pipeline models.PipelineDto) error {

	apiKey, err := config.GetApiKeyForDeployment(deploymentType)
	if err != nil {
		return err
	}

	// Marshal the PipelineDto into YAML
	data, err := yaml.Marshal(pipeline)
	if err != nil {
		return fmt.Errorf("error marshaling pipeline to YAML: %w", err)
	}

	url := getBaseUrl(deploymentType) + "/v1/Pipelines/Commit"
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "text/yaml")
	req.Header.Set("Api-Key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		var respError struct {
			Message string `json:"message"`
		}
		if json.NewDecoder(resp.Body).Decode(&respError) == nil {
			return fmt.Errorf("API error: %s", respError.Message)
		}
		return fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}
