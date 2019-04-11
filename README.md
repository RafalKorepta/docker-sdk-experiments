# docker-sdk-experiments

[![Build Status](https://travis-ci.org/danielpacak/docker-sdk-experiments.svg?branch=master)](https://travis-ci.org/danielpacak/docker-sdk-experiments)

## Apache Kafka integration test

1. Creates a test network docker-sdk
1. Starts zookeeper and kafka containers
2. The kafka container's 9092/tcp port is bound to the random host port
3. The test uses Sarama to connect to kafka container on the random port
4. The test prints Kafka topics
5. The test stops containers

## Kubernetes (K3S) integration test

TODO
