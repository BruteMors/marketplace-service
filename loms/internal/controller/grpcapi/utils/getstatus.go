package utils

import (
	"errors"

	"github.com/BruteMors/marketplace-service/loms/internal/models/order"
	grpcmodels "github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
)

func StringToGRPCOrderStatus(status order.Status) (grpcmodels.OrderStatus, error) {
	switch status {
	case order.OrderStatusNew:
		return grpcmodels.OrderStatus_NEW, nil
	case order.OrderStatusAwaitingPayment:
		return grpcmodels.OrderStatus_AWAITING_PAYMENT, nil
	case order.OrderStatusFailed:
		return grpcmodels.OrderStatus_FAILED, nil
	case order.OrderStatusPayed:
		return grpcmodels.OrderStatus_PAYED, nil
	case order.OrderStatusCancelled:
		return grpcmodels.OrderStatus_CANCELLED, nil
	default:
		return grpcmodels.OrderStatus(0), errors.New("unknown status")
	}
}
