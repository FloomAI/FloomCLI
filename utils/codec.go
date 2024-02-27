package utils

import (
	"FloomCLI/models"
	"gopkg.in/yaml.v3"
	"os"
)

func ParseYaml(yamlFile string) (*models.PipelineDto, error) {
	file, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var config models.PipelineDto
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SerializeYaml(pipeline models.PipelineDto) (string, error) {
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
