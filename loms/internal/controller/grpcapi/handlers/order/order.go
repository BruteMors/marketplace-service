package order

import (
	"context"

	"github.com/BruteMors/marketplace-service/loms/internal/models/order/requests"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/responses"
	"github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
)

type Service interface {
	OrderCreate(ctx context.Context, create *requests.OrderCreate) (orderID int64, err error)
	OrderInfo(ctx context.Context, orderID int64) (responses.OrderInfo, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
}

type GRPCApi struct {
	loms.UnimplementedOrdersServer
	orderService Service
}

func NewOrderGRPCApi(
	orderService Service,
) *GRPCApi {
	return &GRPCApi{
		orderService: orderService,
	}
}
