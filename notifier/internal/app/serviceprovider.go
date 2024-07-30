package app

import (
	"context"
	"log"

	"github.com/BruteMors/marketplace-service/libs/kafka/consumergroup"
	"github.com/BruteMors/marketplace-service/loms/pkg/closer"
	"github.com/BruteMors/marketplace-service/notifier/internal/config"
	"github.com/BruteMors/marketplace-service/notifier/internal/controller/kafka/orderstatus"
	"github.com/BruteMors/marketplace-service/notifier/internal/service/notifier"
)

type serviceProvider struct {
	kafkaConfig             *config.KafkaConfig
	consumerGroup           *consumergroup.ConsumerGroup
	consumerGroupHandler    *consumergroup.Handler
	notifierService         *notifier.Service
	orderStatusKafkaHandler *orderstatus.KafkaHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) KafkaConfig() *config.KafkaConfig {
	if s.kafkaConfig == nil {
		cfg, err := config.NewKafkaConfig()
		if err != nil {
			log.Fatalf("failed to get kafka config: %s", err.Error())
		}

		s.kafkaConfig = cfg
	}

	return s.kafkaConfig
}

func (s *serviceProvider) KafkaConsumerGroup(ctx context.Context) *consumergroup.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := consumergroup.NewConsumerGroup(
			s.KafkaConfig().Brokers(),
			s.KafkaConfig().GroupID(),
			[]string{s.KafkaConfig().GetOrderEventsTopic()},
			s.KafkaConsumerGroupHandler(ctx),
			nil,
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %s", err.Error())
		}

		s.consumerGroup = consumerGroup
	}

	closer.Add(s.consumerGroup.Close)

	return s.consumerGroup
}

func (s *serviceProvider) KafkaConsumerGroupHandler(ctx context.Context) *consumergroup.Handler {
	if s.consumerGroupHandler == nil {
		topicHandlers := make(map[string]consumergroup.TopicHandler)
		topicHandlers[s.KafkaConfig().GetOrderEventsTopic()] = s.OrderStatusKafkaHandler(ctx)

		consumerGroupHandler := consumergroup.NewConsumerGroupHandler(topicHandlers)

		s.consumerGroupHandler = consumerGroupHandler
	}

	return s.consumerGroupHandler
}

func (s *serviceProvider) OrderStatusKafkaHandler(ctx context.Context) *orderstatus.KafkaHandler {
	if s.orderStatusKafkaHandler == nil {
		orderStatusKafkaHandler := orderstatus.NewKafkaHandler(s.NotifierService(ctx))
		s.orderStatusKafkaHandler = orderStatusKafkaHandler
	}

	return s.orderStatusKafkaHandler
}

func (s *serviceProvider) NotifierService(_ context.Context) *notifier.Service {
	if s.notifierService == nil {
		notifierSrv := notifier.NewNotifierService()
		s.notifierService = notifierSrv
	}

	return s.notifierService
}
