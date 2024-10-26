package module

import (
	"encoding/json"
	"errors"
)

type ModuleConfig struct {
	// List of active services
	Modules           []string           `json:"modules"`
	ModuleDefinitions []ModuleDefinition `json:"module_definitions"`
}

type ModuleDefinition struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func ParseServiceConfig(content []byte) (ModuleConfig, error) {
	var config ModuleConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return ModuleConfig{}, err
	}

	for _, service := range config.Modules {
		serviceInDefinitions := false
		for _, definition := range config.ModuleDefinitions {
			if definition.Name == service {
				serviceInDefinitions = true
				break
			}
		}
		if !serviceInDefinitions {
			return config, errors.New("module definition not found: " + service)
		}
	}

	return config, nil
}
