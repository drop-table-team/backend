package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
)

type Module struct {
	ctx    context.Context
	client *client.Client

	definition ModuleDefinition

	container *moduleContainer
}

func (m *Module) IsRunning() bool {
	return m.container != nil
}

func (m *Module) Definition() ModuleDefinition {
	return m.definition
}

func (m *Module) URL() string {
	return fmt.Sprintf("http://%s:%d", m.container.ip, m.definition.Port)
}

func (m *Module) Start() error {
	if m.IsRunning() {
		return errors.New("module is already running")
	}

	var err error
	if m.container, err = initModuleContainer(m.ctx, m.client, m.definition); err != nil {
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

	return nil
}
