package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"log"
	"time"
)

func initModuleContainer(ctx context.Context, client *client.Client, imageRef string, name string) (*moduleContainer, error) {
	// container name is <module name>-<unix timestamp>
	containerName := fmt.Sprintf("%s-%d", name, time.Now().Unix())

	imageSummary, err := client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageRef)),
	})
	if err != nil {
		return nil, err
	}

	// pull image if not found locally
	if len(imageSummary) == 0 {
		out, err := client.ImagePull(ctx, imageRef, image.PullOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("image '%s' not found, pulling it", imageRef)

		// the image is fully downloaded when the `out` reader has reached eof
		var sink [4096]byte
		for {
			_, err = out.Read(sink[:])
			if err != nil && err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
		}

		log.Printf("successfully pulled '%s' image", imageRef)
	}

	response, err := client.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
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
