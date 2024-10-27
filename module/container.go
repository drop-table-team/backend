package module

import (
	"backend/util"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"time"
)

func initModuleContainer(ctx context.Context, client *client.Client, definition ModuleDefinition) (*moduleContainer, error) {
	// container name is <module name>-<unix timestamp>
	containerName := fmt.Sprintf("%s-%d", definition.Name, time.Now().Unix())

	imageSummary, err := client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", definition.Image)),
	})
	if err != nil {
		return nil, err
	}

	// pull image if not found locally
	if len(imageSummary) == 0 {
		out, err := client.ImagePull(ctx, definition.Image, image.PullOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("image '%s' not found, pulling it", definition.Image)

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

		log.Printf("successfully pulled '%s' image", definition.Image)
	}

	envVars := []string{
		fmt.Sprintf("OLLAMA_BASE_URL=%s", util.OllamaUrl),
	}
	var extraHosts []string
	if _, exists := os.LookupEnv("IS_DOCKER"); exists {
		info, err := client.ContainerInspect(ctx, util.UnwrapError(os.Hostname()))
		if err != nil {
			return nil, err
		}
		for _, network := range info.NetworkSettings.Networks {
			if network.NetworkID == networkId {
				envVars = append(envVars, fmt.Sprintf("BACKEND_BASE_URL=http://%s:8080", network.IPAddress))
				break
			}
		}
	} else {
		envVars = append(envVars, fmt.Sprintf("BACKEND_BASE_URL=http://host.docker.internal:8080"))
		extraHosts = append(extraHosts, "host.docker.internal:host-gateway")
	}
	for key, value := range definition.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	response, err := client.ContainerCreate(ctx, &container.Config{
		Image: definition.Image,
		Env:   envVars,
	}, &container.HostConfig{
		ExtraHosts: extraHosts,
		// automatically remove the container when it's stopped
		AutoRemove: true,
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"network": {
				NetworkID: networkId,
			},
		},
	}, nil, containerName)
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
	ip string
}

func (mc *moduleContainer) start() error {
	if mc.id == "" {
		return errors.New("container is already stopped")
	}

	if err := mc.client.ContainerStart(mc.ctx, mc.id, container.StartOptions{}); err != nil {
		return err
	}

	info, err := mc.client.ContainerInspect(mc.ctx, mc.id)
	if err != nil {
		return err
	}
	mc.ip = info.NetworkSettings.Networks["network"].IPAddress

	return nil
}

func (mc *moduleContainer) stop() error {
	if mc.id == "" {
		return errors.New("container is already stopped")
	}

	return mc.client.ContainerStop(mc.ctx, mc.id, container.StopOptions{})
}
