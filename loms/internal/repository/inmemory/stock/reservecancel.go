package stock

import (
	"context"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
)

func (r *Repository) ReserveCancel(ctx context.Context, item []stockmodels.ReserveItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, i := range item {
		stockItem, ok := r.stock[i.SKU]
		if !ok {
			return repository.ErrSKUNotFound
		}

		stockItem.Reserved -= uint64(i.Count)
		r.stock[i.SKU] = stockItem
	}

	return nil
}
