package order

import (
	"context"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
)

type Repository interface {
	Create(ctx context.Context, order ordermodels.NewOrder) (orderID int64, err error)
	SetStatus(ctx context.Context, orderID int64, status ordermodels.Status) error
	GetByID(ctx context.Context, orderID int64) (order ordermodels.Order, err error)
}

type StockService interface {
	Reserve(ctx context.Context, item []ordermodels.Item) error
	ReserveRemove(ctx context.Context, item []stockmodels.ReserveItem) error
	ReserveCancel(ctx context.Context, item []stockmodels.ReserveItem) error
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f func(context.Context) error) error
}

type MQSender interface {
	SendMessage(
		topicName string,
		key []byte,
		message []byte,
		headers map[string]string,
	) (partition int32, offset int64, err error)
}

type StatusOutboxRepository interface {
	CreateOrderStatusChangedEvent(ctx context.Context, orderID int64, status ordermodels.Status) error
	FetchNextOrderStatusChangedEvent(ctx context.Context) (ordermodels.StatusChangedEvent, error)
	MarkOrderStatusChangedEventAsSend(ctx context.Context, eventID int64) error
}

type Service struct {
	orderRepository        Repository
	stockService           StockService
	txManager              TxManager
	mqSender               MQSender
	statusOutboxRepository StatusOutboxRepository
	stopChan               chan struct{}
}

func NewService(
	ctx context.Context,
	repo Repository,
	stockService StockService,
	txManager TxManager,
	mqSender MQSender,
	statusOutboxRepository StatusOutboxRepository,

) *Service {
	s := &Service{
		orderRepository:        repo,
		stockService:           stockService,
		txManager:              txManager,
		mqSender:               mqSender,
		statusOutboxRepository: statusOutboxRepository,
		stopChan:               make(chan struct{}),
	}

	go s.StartStatusChangedEventDispatcher(ctx)

	return s
}

func (s *Service) Close() {
	close(s.stopChan)
}
