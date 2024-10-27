package module

import (
	"backend/util"
	"context"
	"errors"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var networkId string

func NewModuleManager(config ModuleConfig) (ModuleManager, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return ModuleManager{}, err
	}

	var modules []*Module

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

		modules = append(modules, &Module{
			ctx:        ctx,
			client:     cli,
			definition: definition,
		})
	}

	list, err := cli.NetworkList(ctx, network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", util.NetworkName)),
	})
	if err != nil {
		return ModuleManager{}, err
	}

	if len(list) == 0 {
		return ModuleManager{}, errors.New("network not found")
	}

	networkId = list[0].ID

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

	modules []*Module
}

func (m *ModuleManager) FindByName(name string) *Module {
	for _, module := range m.modules {
		if module.definition.Name == name {
			return module
		}
	}
	return nil
}

func (m *ModuleManager) Config() ModuleConfig {
	return m.config
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
