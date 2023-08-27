package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	sleepTimestr := os.Getenv("sleepTime")

	var sleepTime int

	if sleepTimestr == "" {
		sleepTime = 1000
	} else {
		sleepTime, _ = strconv.Atoi(sleepTimestr)
	}

	createTopic(kafkaURL, topic)

	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()
	fmt.Println("start producing ... !!")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced", key)
		}
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func createTopic(kafkaURL, topic string) {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		fmt.Println("cannot connect to kafka")
		panic(err)
	}
	defer conn.Close()

	var numPartitions int
	var replicationFactor int

	numPartitionsStr := os.Getenv("TopicPartitions")
	replicationFactorStr := os.Getenv("ReplicationFactor")

	if numPartitionsStr == "" {
		numPartitions = 1
	} else {
		numPartitions, _ = strconv.Atoi(numPartitionsStr)
	}

	if replicationFactorStr == "" {
		replicationFactor = 1
	} else {
		replicationFactor, _ = strconv.Atoi(replicationFactorStr)
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		fmt.Println("cannot create topic")
		panic(err)
	}
}
