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

func (m *Module) Register() {
	m.registered = true
}

func (m *Module) Unregister() {
	m.registered = false
}

func (m *Module) IsRegistered() bool {
	return m.registered
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
