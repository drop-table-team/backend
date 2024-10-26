package module

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
)

func NewModuleManager(config ModuleConfig) (ModuleManager, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return ModuleManager{}, err
	}

	var modules []Module

	for _, module := range config.Modules {
		var definition ModuleDefinition
		for _, definition = range config.ModuleDefinitions {
			if module == definition.Name {
				break
			}
		}
		if module != definition.Name {
			return ModuleManager{}, errors.New("module definition not found: " + module)
		}

		modules = append(modules, Module{
			ctx:    ctx,
			client: cli,
			name:   definition.Name,
			image:  definition.Image,
		})
	}

	return ModuleManager{
		ctx:     ctx,
		client:  cli,
		config:  config,
		modules: modules,
	}, nil
}

type ModuleManager struct {
	ctx    context.Context
	client *client.Client

	config ModuleConfig

	modules []Module
}

func (m *ModuleManager) FindByName(name string) *Module {
	for _, module := range m.modules {
		if module.name == name {
			return &module
		}
	}
	return nil
}

func (m *ModuleManager) StartAll() error {
	for _, module := range m.modules {
		if err := module.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (m *ModuleManager) StopAll() error {
	for _, module := range m.modules {
		if !module.IsRunning() {
			continue
		}
		if err := module.Stop(); err != nil {
			return err
		}
	}

	return nil
}
