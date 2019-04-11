package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerController struct {
	client *client.Client
	n      *Network
	c      *Container
}

func NewDockerController() (*DockerController, error) {
	client, _ := client.NewEnvClient()

	return &DockerController{
		client: client,
		n: &Network{
			client: client,
		},
		c: &Container{
			client: client,
		},
	}, nil
}

type Container struct {
	client *client.Client
}

func (c *Container) Start(id string) error {
	ctx := context.Background()
	return c.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

func (c *Container) Stop(id string) error {
	ctx := context.Background()
	return c.client.ContainerStop(ctx, id, nil)
}

type ContainerBuilder struct {
	client       *client.Client
	image        string
	name         string
	network      string
	envs         []string
	portBindings map[nat.Port][]nat.PortBinding
}

func (c *Container) Builder() *ContainerBuilder {
	return &ContainerBuilder{
		client: c.client,
	}
}

func (cb *ContainerBuilder) WithName(name string) *ContainerBuilder {
	cb.name = name
	return cb
}

func (cb *ContainerBuilder) WithImage(image string) *ContainerBuilder {
	cb.image = image
	return cb
}

func (cb *ContainerBuilder) WithEnv(name, value string) *ContainerBuilder {
	cb.envs = append(cb.envs, name+"="+value)
	return cb
}

func (cb *ContainerBuilder) WithNetwork(network string) *ContainerBuilder {
	cb.network = network
	return cb
}

func (cb *ContainerBuilder) WithPortBindings(b map[nat.Port][]nat.PortBinding) *ContainerBuilder {
	cb.portBindings = b
	return cb
}

func (cb *ContainerBuilder) Create() (string, error) {
	ctx := context.Background()

	containerCfg := &container.Config{
		Image: cb.image,
		Env:   cb.envs,
	}

	hostCfg := &container.HostConfig{
		AutoRemove:   true,
		NetworkMode:  container.NetworkMode(cb.network),
		PortBindings: cb.portBindings,
	}

	networkdCfg := &network.NetworkingConfig{}

	resp, _ := cb.client.ContainerCreate(ctx, containerCfg, hostCfg, networkdCfg, cb.name)

	return resp.ID, nil
}

func (dc *DockerController) Network() *Network {
	return dc.n
}
func (dc *DockerController) Container() *Container {
	return dc.c
}

type Network struct {
	client *client.Client
}

func (n *Network) Create(name string) (string, error) {
	ctx := context.Background()
	resp, err := n.client.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (n *Network) Remove(name string) error {
	ctx := context.Background()
	return n.client.NetworkRemove(ctx, name)
}
