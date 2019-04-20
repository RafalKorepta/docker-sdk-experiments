package integration

import (
	"fmt"
	"github.com/danielpacak/docker-sdk-experiments/test/common/docker"
	"github.com/danielpacak/docker-sdk-experiments/test/common/kafka"
	"github.com/danielpacak/docker-sdk-experiments/test/common/net"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestKafkaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("This is an integration test")
	}

	// TODO Randomize network name
	const network = "kafka-itest"
	// setup code
	dc, err := docker.NewDockerController()
	require.NoError(t, err)

	_, err = dc.Network().Create(network)
	require.NoError(t, err)

	localIP, err := net.GetLocalIP()
	require.NoError(t, err)

	zookeeperID, err := dc.Container().Builder().
		WithImage("docker.io/confluentinc/cp-zookeeper:5.1.2").
		WithName("zookeeper").
		WithEnv("ZOOKEEPER_CLIENT_PORT", "2181").
		WithEnv("ZOOKEEPER_TICK_TIME", "2000").
		WithEnv("ZOOKEEPER_LOG4J_ROOT_LOGLEVEL", "ERROR").
		WithNetwork(network).
		WithAutoRemove(true).
		Create()
	require.NoError(t, err)

	err = dc.Container().Start(zookeeperID)
	require.NoError(t, err)

	kafkaID, err := dc.Container().Builder().
		WithImage("docker.io/confluentinc/cp-kafka:5.1.2").
		WithName("kafka").
		WithEnvf("KAFKA_ADVERTISED_LISTENERS", "PLAINTEXT://%s:9092", localIP).
		WithEnv("KAFKA_ZOOKEEPER_CONNECT", "zookeeper:2181").
		WithEnv("KAFKA_BROKER_ID", "1").
		WithEnv("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithEnv("KAFKA_AUTO_CREATE_TOPICS_ENABLE", "false").
		WithNetwork(network).
		WithPortBindings(map[nat.Port][]nat.PortBinding{
			"9092/tcp": {
				nat.PortBinding{HostPort: "9092"}},
		}, ).
		WithAutoRemove(true).
		Create()
	require.NoError(t, err)

	err = dc.Container().Start(kafkaID)
	require.NoError(t, err)

	t.Run("Should create and list topics", func(t *testing.T) {
		brokerAddr := fmt.Sprintf("%s:9092", localIP)

		time.Sleep(15 * time.Second)

		admin, err := kafka.NewAdmin(t, []string{brokerAddr})
		require.NoError(t, err)

		err = admin.CreateTopic("test.topic.1")
		require.NoError(t, err)

		err = admin.CreateTopic("test.topic.2")
		require.NoError(t, err)

		time.Sleep(3 * time.Second)

		topics, err := admin.GetTopicNames()
		require.NoError(t, err)

		t.Logf("Kafka topics: %v", topics)
		assert.Contains(t, topics, "test.topic.1")
		assert.Contains(t, topics, "test.topic.2")

		err = admin.Close()
		require.NoError(t, err)
	})

	// tear-down code
	err = dc.Container().Stop(kafkaID)
	require.NoError(t, err)

	err = dc.Container().Stop(zookeeperID)
	require.NoError(t, err)

	err = dc.Network().Remove(network)
	require.NoError(t, err)

}
