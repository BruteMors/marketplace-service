package app

import (
	"context"
	"log"

	"github.com/BruteMors/marketplace-service/libs/kafka"
	"github.com/BruteMors/marketplace-service/libs/kafka/producer"
	"github.com/BruteMors/marketplace-service/loms/internal/config"
	"github.com/BruteMors/marketplace-service/loms/internal/controller/grpcapi/handlers/order"
	"github.com/BruteMors/marketplace-service/loms/internal/controller/grpcapi/handlers/stock"
	inMemoryorderRepository "github.com/BruteMors/marketplace-service/loms/internal/repository/inmemory/order"
	inMemorystockRepository "github.com/BruteMors/marketplace-service/loms/internal/repository/inmemory/stock"
	orderRepository "github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/outbox"
	stockRepository "github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock"
	orderService "github.com/BruteMors/marketplace-service/loms/internal/service/order"
	stockService "github.com/BruteMors/marketplace-service/loms/internal/service/stock"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"github.com/BruteMors/marketplace-service/loms/pkg/closer"
)

type serviceProvider struct {
	grpcConfig              *config.GRPCServerConfig
	httpConfig              *config.HTTPServerConfig
	kafkaConfig             *config.KafkaConfig
	mqSyncProducer          *producer.SyncProducer
	stockGrpcApi            *stock.GRPCApi
	stockService            *stockService.Service
	inMemoryStockRepository *inMemorystockRepository.Repository
	stockRepository         *stockRepository.Repository
	orderGrpcApi            *order.GRPCApi
	orderService            *orderService.Service
	inMemoryOrderRepository *inMemoryorderRepository.Repository
	orderRepository         *orderRepository.Repository
	pgConfig                config.PGConfig
	dbClient                *pg.Client
	txManager               transaction.TxManager
	outboxRepository        *outbox.Repository
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCServerConfig() *config.GRPCServerConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCServerConfig()
		if err != nil {
			log.Fatalf("failed to get gRPC config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPServerConfig() *config.HTTPServerConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPServerConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
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

func (s *serviceProvider) KafkaSyncProducer(_ context.Context) *producer.SyncProducer {
	if s.mqSyncProducer == nil {
		syncProducer, err := producer.NewSyncProducer(kafka.Config{Brokers: s.KafkaConfig().Brokers()}, nil)
		if err != nil {
			log.Fatalf("failed to create kafka producer: %s", err.Error())
		}

		s.mqSyncProducer = syncProducer
	}

	closer.Add(s.mqSyncProducer.Close)

	return s.mqSyncProducer
}

func (s *serviceProvider) OrderRepository(ctx context.Context) *orderRepository.Repository {
	if s.orderRepository == nil {
		orderRepo := orderRepository.NewRepository(s.DBClient(ctx))
		s.orderRepository = orderRepo
	}

	return s.orderRepository
}

func (s *serviceProvider) InMemoryOrderRepository(_ context.Context) *inMemoryorderRepository.Repository {
	if s.inMemoryOrderRepository == nil {
		orderRepo, err := inMemoryorderRepository.NewRepository()
		if err != nil {
			log.Fatalf("failed to get order repository: %s", err.Error())
		}

		s.inMemoryOrderRepository = orderRepo
	}

	return s.inMemoryOrderRepository
}

func (s *serviceProvider) OrderService(ctx context.Context) *orderService.Service {
	if s.orderService == nil {
		orderSrv := orderService.NewService(
			ctx,
			s.OrderRepository(ctx),
			s.StockService(ctx),
			s.TxManager(ctx),
			s.KafkaSyncProducer(ctx),
			s.OutboxRepository(ctx),
		)

		s.orderService = orderSrv

		closer.Add(func() error {
			orderSrv.Close()
			return nil
		})
	}

	return s.orderService
}

func (s *serviceProvider) OrderGRPCApi(ctx context.Context) *order.GRPCApi {
	if s.orderGrpcApi == nil {
		s.orderGrpcApi = order.NewOrderGRPCApi(
			s.OrderService(ctx),
		)
	}

	return s.orderGrpcApi
}

func (s *serviceProvider) StockRepository(ctx context.Context) *stockRepository.Repository {
	if s.stockRepository == nil {
		stockRepo := stockRepository.NewRepository(s.DBClient(ctx))
		s.stockRepository = stockRepo
	}

	return s.stockRepository
}

func (s *serviceProvider) InMemoryStockRepository(_ context.Context) *inMemorystockRepository.Repository {
	if s.inMemoryStockRepository == nil {
		stockRepo, err := inMemorystockRepository.NewRepository()
		if err != nil {
			log.Fatalf("failed to get stock repository: %s", err.Error())
		}

		s.inMemoryStockRepository = stockRepo
	}

	return s.inMemoryStockRepository
}

func (s *serviceProvider) StockService(ctx context.Context) *stockService.Service {
	if s.stockService == nil {
		stockSvc := stockService.NewService(
			s.StockRepository(ctx),
			s.TxManager(ctx),
		)

		s.stockService = stockSvc
	}

	return s.stockService
}

func (s *serviceProvider) StockGRPCApi(ctx context.Context) *stock.GRPCApi {
	if s.stockGrpcApi == nil {
		s.stockGrpcApi = stock.NewStockGRPCApi(
			s.StockService(ctx),
		)
	}

	return s.stockGrpcApi
}

func (s *serviceProvider) OutboxRepository(ctx context.Context) *outbox.Repository {
	if s.outboxRepository == nil {
		outboxRepo := outbox.NewRepository(s.DBClient(ctx))
		s.outboxRepository = outboxRepo
	}

	return s.outboxRepository
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) *pg.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().MasterDSN(), s.PGConfig().ReplicaDSNs())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.MasterDB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error (master): %s", err.Error())
		}

		for _, replicaDBC := range cl.ReplicaDBs() {
			err = replicaDBC.Ping(ctx)
			if err != nil {
				log.Fatalf("ping error (replica): %s", err.Error())
			}
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) transaction.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx))
	}

	return s.txManager
}
