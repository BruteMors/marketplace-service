package config

import (
	"errors"
	"os"
)

const (
	kafkaBrokersEnvName     = "KAFKA_BROKERS"
	orderEventsTopicEnvName = "ORDER_EVENTS_TOPIC"
	groupIDEnvName          = "GROUP_ID"
)

type KafkaConfig struct {
	brokers          []string
	orderEventsTopic string
	groupID          string
}

func NewKafkaConfig() (*KafkaConfig, error) {
	brokers := os.Getenv(kafkaBrokersEnvName)
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers not found")
	}

	orderEventsTopic := os.Getenv(orderEventsTopicEnvName)
	if len(orderEventsTopic) == 0 {
		return nil, errors.New("order events topic not found")
	}

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("group id not found")
	}

	return &KafkaConfig{
		brokers:          []string{brokers},
		orderEventsTopic: orderEventsTopic,
		groupID:          groupID,
	}, nil
}

func (c *KafkaConfig) Brokers() []string {
	return c.brokers
}

func (c *KafkaConfig) GetOrderEventsTopic() string {
	return c.orderEventsTopic
}

func (c *KafkaConfig) GroupID() string {
	return c.groupID
}
