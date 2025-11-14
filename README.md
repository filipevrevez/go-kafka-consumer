# GO Kafka Consumer for KEDA Demo

This is a simple Go application that consumes messages from a Kafka topic and prints them to the console. It is intended to be used as a demo for [KEDA](https://keda.sh/).

## Prerequisites

- Kubernetes Cluster (tested on [KIND](https://kind.sigs.k8s.io/))
- [KEDA](https://keda.sh/) installed in the cluster
- [Kafka](https://kafka.apache.org/) installed in the cluster
- Internet connection to pull the Docker image

## Usage

Launch the kind cluster:

```bash
kind create cluster --config kind-config.yaml
```

Install KEDA:

```bash
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda --namespace keda --create-namespace
```

Next, install Kafka. I used the [Ricardo Aires Helm Chart](https://ricardo-aires.github.io/helm-charts/charts/kafka/):

```bash
helm repo add rhcharts https://ricardo-aires.github.io/helm-charts/
helm upgrade --install kafka rhcharts/kafka -f kafka-values.yaml --namespace kafka --create-namespace
```

Using [Conduktor](https://www.conduktor.io/), create a topic called `test-topic` with 3 partitions and 3 replicas.

Next, put some data via Conduktor into the topic. You should see the consumer print the messages to the console.

Finally, deploy the consumer to the cluster:

```bash
kubectl apply -f deployment.yaml
```

And scale the deployment to 0:

```bash
kubectl scale deployment kafka-consumer --replicas=0
```

After that, lets launch the KEDA ScaledObject:

```bash
kubectl apply -f keda.yaml
```

Now, if you put some data into the topic, you should see the consumer print the messages to the console again.
You should also see the ScaledObject scale the deployment to some replicas depending on your configuration.

## Repositories used

[Segmentio kakfka go](https://github.com/segmentio/kafka-go)
