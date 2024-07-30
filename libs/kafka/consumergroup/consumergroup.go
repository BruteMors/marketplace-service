package consumergroup

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (c *ConsumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
				slog.ErrorContext(ctx, "Error from consume: %v\n", err)
			}
			if ctx.Err() != nil {
				slog.InfoContext(ctx, "[consumer-group]: ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()
}

func NewConsumerGroup(
	brokers []string,
	groupID string,
	topics []string,
	consumerGroupHandler sarama.ConsumerGroupHandler,
	opts ...Option,
) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = false

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		err := opt.Apply(config)
		if err != nil {
			return nil, err
		}
	}

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}

func (c *ConsumerGroup) Close() error {
	return c.ConsumerGroup.Close()
}
