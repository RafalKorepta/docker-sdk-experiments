# docker-sdk-experiments

[![Build Status](https://travis-ci.org/danielpacak/docker-sdk-experiments.svg?branch=master)](https://travis-ci.org/danielpacak/docker-sdk-experiments)

## Apache Kafka integration test - custom bridge network

1. Create the `kafka-itest` network.
1. Start `zookeeper` and `kafka` Docker containers.
2. Bind `kafka` container's port `9092/tcp` to a random host port.
3. Use [Sarama](https://github.com/Shopify/sarama) to connect to the `kafka`
   container on the random port.
4. Print Kafka topics.
5. Stop `zookeeper` and `kafka` containers.
6. Remove `kafka-itest` network.

```text
=== RUN   TestKafkaIntegration
--- PASS: TestKafkaIntegration (23.60s)
=== RUN   TestKafkaIntegration/Should_list_default_topics
    --- PASS: TestKafkaIntegration/Should_list_default_topics (15.01s)
        kafka_integration_test.go:69: Kafka topics: [__confluent.support.metrics]
PASS

Process finished with exit code 0
```

## Apache Kafka integration test - host network

> NB This won't work on macOS

1. Start `zookeeper` and `kafka` Docker containers and attach them to the `host`
   network.
2. Use Sarama to connect to Kafka at `localhost:9092`.
3. Print Kafka topics.
4. Stop `kafka` and `zookeeper` containers.

## Kubernetes (K3S) integration test

1. Create the `k8s-itest` network.
2. Create and start [`k3s`](https://github.com/rancher/k3s) Docker container.
3. Connect to K3S API with K8S client-go.
4. List all namespaces.
5. Stop K3S container.
6. Remove the `k3s-itest` network.

```text
=== RUN   TestK8SIntegration
--- PASS: TestK8SIntegration (5.20s)
=== RUN   TestK8SIntegration/Should_list_namespaces
    --- PASS: TestK8SIntegration/Should_list_namespaces (3.52s)
        k8s_integration_test.go:55: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        io.go:33: Waiting for file /tmp/kubeconfig.yaml to become available
        k8s_integration_test.go:55: Done waiting for file /tmp/kubeconfig.yaml to become available
        k8s_integration_test.go:68: namespace: default, status: {Active}
        k8s_integration_test.go:68: namespace: kube-public, status: {Active}
        k8s_integration_test.go:68: namespace: kube-system, status: {Active}
PASS

Process finished with exit code 0
```

## TODO

1. Add example with accessing files from named volume.

## Read

1. [Networking features in Docker Desktop for Mac](https://docs.docker.com/docker-for-mac/networking/)
2. [Arquillian Cube](http://arquillian.org/arquillian-cube/)
