/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
package events

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

// GetKafkaConfig returns a Kafka configuration from environment variables
func GetKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Retry.Max = 3
	return config
}

// GetKafkaBrokers returns the list of Kafka brokers from environment
func GetKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	return strings.Split(brokers, ",")
}

// GetKafkaTopic returns the Kafka topic name from environment
func GetKafkaTopic() string {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "mushop-events"
	}
	return topic
}
