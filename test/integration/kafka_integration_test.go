package integration

import (
	"github.com/Shopify/sarama"
	"github.com/danielpacak/docker-sdk-experiments/net"
	"github.com/danielpacak/docker-sdk-experiments/test/common/docker"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestKafkaIntegration(t *testing.T) {
	// setup code
	dc, err := docker.NewDockerController()
	require.NoError(t, err)

	// TODO Randomize network name
	_, err = dc.Network().Create("docker-sdk")
	require.NoError(t, err)

	zookeeperID, err := dc.Container().Builder().
		WithImage("docker.io/confluentinc/cp-zookeeper:5.1.2").
		WithName("zookeeper").
		WithEnv("ZOOKEEPER_CLIENT_PORT", "2181").
		WithEnv("ZOOKEEPER_TICK_TIME", "2000").
		WithEnv("ZOOKEEPER_LOG4J_ROOT_LOGLEVEL", "ERROR").
		WithNetwork("docker-sdk").
		Create()
	require.NoError(t, err)

	err = dc.Container().Start(zookeeperID)
	require.NoError(t, err)

	port, err := net.GetFreePort()
	require.NoError(t, err)

	kafkaID, err := dc.Container().Builder().
		WithImage("docker.io/confluentinc/cp-kafka:5.1.2").
		WithName("kafka").
		WithEnv("KAFKA_ADVERTISED_LISTENERS", "PLAINTEXT://kafka:9092").
		WithEnv("KAFKA_ZOOKEEPER_CONNECT", "zookeeper:2181").
		WithEnv("KAFKA_BROKER_ID", "1").
		WithEnv("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithNetwork("docker-sdk").
		WithPortBindings(map[nat.Port][]nat.PortBinding{
			"9092/tcp": {
				nat.PortBinding{HostPort: strconv.Itoa(port)}},
		}, ).
		Create()
	require.NoError(t, err)

	err = dc.Container().Start(kafkaID)
	require.NoError(t, err)

	t.Run("Should list default topics", func(t *testing.T) {
		// TODO Check that Kafka listens on :9092
		time.Sleep(15 * time.Second)
		config := sarama.NewConfig()

		saramaClient, err := sarama.NewClient([]string{"localhost:" + strconv.Itoa(port)}, config)
		require.NoError(t, err)

		topics, err := saramaClient.Topics()
		require.NoError(t, err)

		t.Logf("Kafka topics: %v", topics)

		err = saramaClient.Close()
		require.NoError(t, err)
	})

	// tear-down code
	err = dc.Container().Stop(kafkaID)
	require.NoError(t, err)
	err = dc.Container().Stop(zookeeperID)
	require.NoError(t, err)

	err = dc.Network().Remove("docker-sdk")
	require.NoError(t, err)

}
