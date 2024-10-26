package module

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
)

type Module struct {
	ctx    context.Context
	client *client.Client

	name  string
	image string

	container *moduleContainer

	registered bool
}

func (m *Module) IsRunning() bool {
	return m.container != nil
}

func (m *Module) Start() error {
	if m.IsRunning() {
		return errors.New("module is already running")
	}

	var err error
	if m.container, err = initModuleContainer(m.ctx, m.client, m.image, m.name); err != nil {
		return err
	}

	return m.container.start()
}

func (m *Module) Stop() error {
	if !m.IsRunning() {
		return errors.New("module container has not been started")
	}

	if err := m.container.stop(); err != nil {
		return err
	}

	m.container = nil
	m.registered = false

	return nil
}

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

func (s *ModuleManager) StartAll() error {
	for _, module := range s.modules {
		if err := module.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (s *ModuleManager) StopAll() error {
	for _, module := range s.modules {
		if !module.IsRunning() {
			continue
		}
		if err := module.Stop(); err != nil {
			return err
		}
	}

	return nil
}
