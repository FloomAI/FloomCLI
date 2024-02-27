package models

// PipelineDto represents the structure of a pipeline.
type PipelineDto struct {
	Kind     string             `yaml:"kind"`
	Pipeline PipelineDetailsDto `yaml:"pipeline"`
}

type PipelineDetailsDto struct {
	Name     string                   `yaml:"name"`
	Model    []PluginConfigurationDto `yaml:"model,omitempty"`
	Prompt   *PromptStageDto          `yaml:"prompt,omitempty"`
	Response *ResponseStageDto        `yaml:"response,omitempty"`
	Global   []PluginConfigurationDto `yaml:"global,omitempty"`
}

type PromptStageDto struct {
	Template     *PluginConfigurationDto  `yaml:"template,omitempty"`
	Context      []PluginConfigurationDto `yaml:"context,omitempty"`
	Optimization []PluginConfigurationDto `yaml:"optimization,omitempty"`
	Validation   []PluginConfigurationDto `yaml:"validation,omitempty"`
}

type ResponseStageDto struct {
	Format     []PluginConfigurationDto `yaml:"format,omitempty"`
	Validation []PluginConfigurationDto `yaml:"validation,omitempty"`
}

type PluginConfigurationDto struct {
	Package       string                 `yaml:"package"`
	Configuration map[string]interface{} `yaml:"configuration"`
}

func (p *PluginConfigurationDto) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Temporary structure to capture all fields
	var temp struct {
		Package string                 `yaml:"package"`
		Other   map[string]interface{} `yaml:",inline"`
	}

	if err := unmarshal(&temp); err != nil {
		return err
	}

	p.Package = temp.Package
	p.Configuration = make(map[string]interface{})

	// Move all fields except 'package' to Configuration
	for key, value := range temp.Other {
		p.Configuration[key] = value
	}

	return nil
}

func (p PluginConfigurationDto) MarshalYAML() (interface{}, error) {
	// Start with copying the Configuration map
	output := make(map[string]interface{}, len(p.Configuration)+1)

	// Add the package field to the output map
	output["package"] = p.Package

	// Add all other Configuration fields to the output map
	for key, value := range p.Configuration {
		output[key] = value
	}

	return output, nil
}
