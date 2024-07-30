package orderstatus

import (
	"context"

	"github.com/BruteMors/marketplace-service/notifier/internal/models/order"
)

type Service interface {
	ProcessOrderStatus(ctx context.Context, event order.StatusChangedEvent) error
}

type KafkaHandler struct {
	orderService Service
}

func NewKafkaHandler(
	orderService Service,
) *KafkaHandler {
	return &KafkaHandler{
		orderService: orderService,
	}
}
