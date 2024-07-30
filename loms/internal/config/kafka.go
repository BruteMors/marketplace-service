package config

import (
	"errors"
	"os"
)

const (
	kafkaBrokersEnvName     = "KAFKA_BROKERS"
	orderEventsTopicEnvName = "ORDER_EVENTS_TOPIC"
)

type KafkaConfig struct {
	brokers []string
}

func NewKafkaConfig() (*KafkaConfig, error) {
	brokers := os.Getenv(kafkaBrokersEnvName)
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers not found")
	}

	return &KafkaConfig{
		brokers: []string{brokers},
	}, nil
}

func (c *KafkaConfig) Brokers() []string {
	return c.brokers
}

func GetSendStatusChangedEventTopic() string {
	return os.Getenv(orderEventsTopicEnvName)
}
