package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"io"
	"os"
)

const (
	zookeeperImage = "docker.io/confluentinc/cp-zookeeper:5.1.2"
)

func main() {
	ctx := context.Background()
	docker, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	reader, err := docker.ImagePull(ctx, zookeeperImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: zookeeperImage,
		Env:   []string{"ZOOKEEPER_CLIENT_PORT=2181"},
	}, &container.HostConfig{NetworkMode: "host"}, nil, "zookeeper")
	if err != nil {
		panic(err)
	}

	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err := docker.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow:true})
	if err != nil {
		panic(err)
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		panic(err)
	}

}
