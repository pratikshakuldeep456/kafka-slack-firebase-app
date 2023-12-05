package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	//"github.com/Shopify/sarama"
)

const (
	kafkaBrokers   = "your_kafka_brokers"
	kafkaGroupID   = "your_group_id"
	kafkaTopic     = "your_kafka_topic"
	slackWebhook   = "your_slack_webhook_url"
	firebaseConfig = "your_firebase_config_file.json"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{kafkaBrokers}, config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(kafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("Error consuming Kafka partition: %v\n", err)
		return
	}
	defer partitionConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go kafkaMessageProcessor(ctx, partitionConsumer, wg)

	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case sig := <-signals:
			fmt.Printf("Received signal %v, shutting down...\n", sig)
			cancel()
		}
	}()

	wg.Wait()
}

func kafkaMessageProcessor(ctx context.Context, consumer sarama.PartitionConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-consumer.Messages():
			processMessage(msg.Value)
		}
	}
}

func processMessage(message []byte) {
	payload := string(message)
	fmt.Printf("Received Kafka message: %s\n", payload)

	sendSlackNotification(payload)
	sendFirebaseNotification(payload)
}

func sendSlackNotification(message string) {
	// Implement Slack notification logic here
	fmt.Printf("Sending Slack notification: %s\n", message)
}

func sendFirebaseNotification(message string) {
	// Implement Firebase notification logic here
	fmt.Printf("Sending Firebase notification: %s\n", message)
}
