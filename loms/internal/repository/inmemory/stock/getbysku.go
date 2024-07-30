package stock

import (
	"context"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
)

func (r *Repository) GetBySKU(ctx context.Context, skuID uint32) (stockmodels.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.stock[skuID]
	if !ok {
		return stockmodels.Item{}, repository.ErrSKUNotFound
	}

	return stockmodels.Item{
		SKU:        item.SKU,
		TotalCount: item.TotalCount,
		Reserved:   item.Reserved,
	}, nil
}
