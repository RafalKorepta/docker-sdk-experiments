package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/danielpacak/docker-sdk-experiments/net"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	zookeeperImage = "docker.io/confluentinc/cp-zookeeper:5.1.2"
	kafkaImage     = "docker.io/confluentinc/cp-kafka:5.1.2"
)

func main() {
	createNetwork()
	startZookeeperContainer()

	port, err := net.GetFreePort()

	startKafkaContainer(port)

	time.Sleep(15 * time.Second)
	config := sarama.NewConfig()

	saramaClient, err := sarama.NewClient([]string{"localhost:" + strconv.Itoa(port)}, config)
	if err != nil {
		panic(err)
	}

	topics, err := saramaClient.Topics()
	if err != nil {
		panic(err)
	}
	log.Printf("Kafka topics: %v", topics)

	err = saramaClient.Close()
	if err != nil {
		panic(err)
	}

	//removeNetwork()
}

func createNetwork() {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	_, err = docker.NetworkCreate(ctx, "docker-sdk", types.NetworkCreate{})
	if err != nil {
		panic(err)
	}
}

func removeNetwork() {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	err = docker.NetworkRemove(ctx, "docker-sdk")
	if err != nil {
		panic(err)
	}
}

func startKafkaContainer(port int) {
	ctx := context.Background()
	docker, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	reader, err := docker.ImagePull(ctx, kafkaImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	hostCfg := &container.HostConfig{
		AutoRemove:  false,
		NetworkMode: "docker-sdk",
		PortBindings: map[nat.Port][]nat.PortBinding{
			"9092/tcp": {
				nat.PortBinding{HostPort: strconv.Itoa(port)}},
		},
	}
	networkdCfg := &network.NetworkingConfig{
	}

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: kafkaImage,
		Env: []string{
			"KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092",
			"KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181",
			"KAFKA_BROKER_ID=1",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
		},
	}, hostCfg, networkdCfg, "kafka")
	if err != nil {
		panic(err)
	}

	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func startZookeeperContainer() {
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

	hostCfg := &container.HostConfig{
		AutoRemove:  false,
		NetworkMode: "docker-sdk",
	}
	networkdCfg := &network.NetworkingConfig{

	}

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: zookeeperImage,
		Env: []string{
			"ZOOKEEPER_CLIENT_PORT=2181",
			"ZOOKEEPER_TICK_TIME=2000",
			"ZOOKEEPER_LOG4J_ROOT_LOGLEVEL=ERROR",
		},
	}, hostCfg, networkdCfg, "zookeeper")
	if err != nil {
		panic(err)
	}

	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}
