package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"time"
)

func initModuleContainer(ctx context.Context, client *client.Client, image string, name string) (*moduleContainer, error) {
	// container name is <module name>-<unix timestamp>
	containerName := fmt.Sprintf("%s-%d", name, time.Now().Unix())

	response, err := client.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, &container.HostConfig{
		// automatically remove the container when it's stopped
		AutoRemove: true,
	}, nil, nil, containerName)
	if err != nil {
		return nil, err
	}

	return &moduleContainer{
		ctx:    ctx,
		client: client,
		id:     response.ID,
	}, nil
}

type moduleContainer struct {
	ctx    context.Context
	client *client.Client

	id string
}

func (mc *moduleContainer) start() error {
	if mc.id == "" {
		return errors.New("container is already stopped")
	}

	return mc.client.ContainerStart(mc.ctx, mc.id, container.StartOptions{})
}

func (mc *moduleContainer) stop() error {
	if mc.id == "" {
		return errors.New("container is already stopped")
	}

	return mc.client.ContainerStop(mc.ctx, mc.id, container.StopOptions{})
}
