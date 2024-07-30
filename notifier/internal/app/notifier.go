package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/BruteMors/marketplace-service/libs/kafka/consumergroup"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/notifier/internal/config"
	"github.com/BruteMors/marketplace-service/notifier/pkg/closer"
)

type NotifierApp struct {
	serviceProvider *serviceProvider
	consumerGroup   *consumergroup.ConsumerGroup
	shutdownTracer  func(context.Context) error
	wg              sync.WaitGroup
}

func NewNotifier(ctx context.Context) (*NotifierApp, error) {
	a := &NotifierApp{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *NotifierApp) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.shutdownTracer(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	c.consumerGroup.Run(context.Background(), &wg)

	wg.Wait()

	return nil
}

func (c *NotifierApp) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		c.initConfig,
		c.initServiceProvider,
		c.initTracing,
		c.initKafkaConsumerGroup,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *NotifierApp) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (c *NotifierApp) initServiceProvider(_ context.Context) error {
	c.serviceProvider = newServiceProvider()
	return nil
}

func (c *NotifierApp) initKafkaConsumerGroup(ctx context.Context) error {
	c.consumerGroup = c.serviceProvider.KafkaConsumerGroup(ctx)
	return nil
}

func (c *NotifierApp) initTracing(_ context.Context) error {
	c.shutdownTracer = tracing.InitTracer()
	return nil
}
