package kafka

import (
	"github.com/Shopify/sarama"
	"testing"
)

type Admin struct {
	t            *testing.T
	clusterAdmin sarama.ClusterAdmin
}

func NewAdmin(t *testing.T, addrs []string) (*Admin, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_0_0
	cfg.ClientID = "test-admin"
	clusterAdmin, err := sarama.NewClusterAdmin(addrs, cfg)
	if err != nil {
		return nil, err
	}
	return &Admin{
		t:            t,
		clusterAdmin: clusterAdmin,
	}, nil
}

func (a *Admin) Close() error {
	return a.clusterAdmin.Close()
}

func (a *Admin) CreateTopic(name string) error {
	return a.clusterAdmin.CreateTopic(name, &sarama.TopicDetail{NumPartitions: 1, ReplicationFactor: 1}, false)
}

func (a *Admin) GetTopicNames() ([]string, error) {
	var names []string
	topics, err := a.clusterAdmin.ListTopics()
	if err != nil {
		return nil, err
	}
	for k := range topics {
		names = append(names, k)
	}
	return names, nil
}
