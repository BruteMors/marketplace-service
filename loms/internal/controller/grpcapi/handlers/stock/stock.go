package stock

import (
	"context"

	"github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
)

type Service interface {
	StocksInfo(ctx context.Context, sku uint32) (count uint64, err error)
}

type GRPCApi struct {
	loms.UnimplementedStockServer
	stockService Service
}

func NewStockGRPCApi(
	stockService Service,
) *GRPCApi {
	return &GRPCApi{
		stockService: stockService,
	}
}
