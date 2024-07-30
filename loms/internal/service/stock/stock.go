package stock

import (
	"context"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
)

type Repository interface {
	GetBySKU(ctx context.Context, skuID uint32) (stockmodels.Item, error)
	Reserve(ctx context.Context, item []stockmodels.ReserveItem) error
	ReserveRemove(ctx context.Context, item []stockmodels.ReserveItem) error
	ReserveCancel(ctx context.Context, item []stockmodels.ReserveItem) error
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f func(context.Context) error) error
}

type Service struct {
	stockRepository Repository
	txManager       TxManager
}

func NewService(
	repo Repository,
	txManager TxManager,
) *Service {
	return &Service{
		stockRepository: repo,
		txManager:       txManager,
	}
}
