package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"io"
	"os"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"
)

const (
	volumeName = "test_volume"
	imageRef   = "docker.io/rancher/k3s:v0.2.0"
)

func TestK8SIntegration(t *testing.T) {
	ctx := context.Background()
	docker, err := client.NewEnvClient()
	require.NoError(t, err)

	t.Logf("Creating volume %s ...", volumeName)
	_, err = docker.VolumeCreate(ctx, volume.VolumesCreateBody{Name: volumeName})
	require.NoError(t, err)
	t.Logf("Done creating volume %s", volumeName)

	reader, err := docker.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	require.NoError(t, err)
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	require.NoError(t, err)

	t.Logf("Creating container ...")
	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		Env: [] string{
			"K3S_CLUSTER_SECRET=somethingtotallyrandom",
			"K3S_KUBECONFIG_OUTPUT=/tmp/kubeconfig.yaml",
			"K3S_KUBECONFIG_MODE=666",
		},
		Cmd: []string{"server", "--disable-agent", "--https-listen-port", "60443"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: "/tmp",
			},
		},
	}, &network.NetworkingConfig{}, "")
	require.NoError(t, err)
	t.Logf("Done creating container: %s", resp.ID)

	t.Logf("Starting container %s ...", resp.ID)
	err = docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	require.NoError(t, err)
	t.Logf("Done starting container %s", resp.ID)

	time.Sleep(30 * time.Second)

	t.Logf("Stopping container %s ...", resp.ID)
	stopTimout := 30 * time.Second
	err = docker.ContainerStop(ctx, resp.ID, &stopTimout)
	require.NoError(t, err)
	t.Logf("Done stopping container %s", resp.ID)

	t.Logf("Removing container %s ...", resp.ID)
	err = docker.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{RemoveVolumes: true})
	require.NoError(t, err)
	t.Logf("Done removing container %s", resp.ID)

	t.Logf("Removing volume %s ...", volumeName)
	err = docker.VolumeRemove(ctx, volumeName, false)
	require.NoError(t, err)
	t.Logf("Done removing volume %s", volumeName)
}
