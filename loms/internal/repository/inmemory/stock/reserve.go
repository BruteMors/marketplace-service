package stock

import (
	"context"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
)

func (r *Repository) Reserve(ctx context.Context, items []stockmodels.ReserveItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		stockItem, ok := r.stock[item.SKU]
		if !ok {
			return repository.ErrSKUNotFound
		}

		if stockItem.Reserved+uint64(item.Count) > stockItem.TotalCount {
			return repository.ErrInsufficientStock
		}

		stockItem.Reserved += uint64(item.Count)
		r.stock[item.SKU] = stockItem
	}

	return nil
}
