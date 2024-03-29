package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
)

type DockerController struct {
	client *client.Client
	n      *Network
	image  *ImageSvc
	c      *Container
}

func NewDockerController() (*DockerController, error) {
	client, _ := client.NewEnvClient()

	return &DockerController{
		client: client,
		n: &Network{
			client: client,
		},
		image: &ImageSvc{
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
	cmd          []string
	name         string
	network      string
	envs         []string
	exposedPorts nat.PortSet
	portBindings map[nat.Port][]nat.PortBinding
	mounts       []mount.Mount
	autoRemove   bool
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

func (cb *ContainerBuilder) WithCmd(cmd []string) *ContainerBuilder {
	cb.cmd = cmd
	return cb
}

func (cb *ContainerBuilder) WithEnv(name, value string) *ContainerBuilder {
	cb.envs = append(cb.envs, name+"="+value)
	return cb
}

func (cb *ContainerBuilder) WithEnvf(name, format string, a ...interface{}) *ContainerBuilder {
	cb.envs = append(cb.envs, name+"="+fmt.Sprintf(format, a))
	return cb
}

func (cb *ContainerBuilder) WithNetwork(network string) *ContainerBuilder {
	cb.network = network
	return cb
}

func (cb *ContainerBuilder) WithExposedPorts(exposedPorts nat.PortSet) *ContainerBuilder {
	cb.exposedPorts = exposedPorts
	return cb
}

func (cb *ContainerBuilder) WithPortBindings(b map[nat.Port][]nat.PortBinding) *ContainerBuilder {
	cb.portBindings = b
	return cb
}

func (cb *ContainerBuilder) WithMounts(mounts []mount.Mount) *ContainerBuilder {
	cb.mounts = mounts
	return cb
}

func (cb *ContainerBuilder) WithAutoRemove(autoRemove bool) *ContainerBuilder {
	cb.autoRemove = autoRemove
	return cb
}

func (cb *ContainerBuilder) Create() (string, error) {
	ctx := context.Background()

	containerCfg := &container.Config{
		Image:        cb.image,
		Env:          cb.envs,
		Cmd:          cb.cmd,
		ExposedPorts: cb.exposedPorts,
	}

	hostCfg := &container.HostConfig{
		AutoRemove:   cb.autoRemove,
		NetworkMode:  container.NetworkMode(cb.network),
		PortBindings: cb.portBindings,
		Mounts:       cb.mounts,
	}

	networkdCfg := &network.NetworkingConfig{}

	resp, _ := cb.client.ContainerCreate(ctx, containerCfg, hostCfg, networkdCfg, cb.name)

	return resp.ID, nil
}

func (dc *DockerController) Network() *Network {
	return dc.n
}

func (dc *DockerController) Image() *ImageSvc {
	return dc.image
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

// Remove removes a given network.
func (n *Network) Remove(name string) error {
	ctx := context.Background()
	return n.client.NetworkRemove(ctx, name)
}

type ImageSvc struct {
	client *client.Client
}

// Pull pulls an image or a repository from a registry.
func (n *ImageSvc) Pull(name string) error {
	ctx := context.Background()
	out, err := n.client.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, out)
	return err
}
