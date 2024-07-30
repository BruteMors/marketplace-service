package order

import (
	"context"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status ordermodels.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[uint64(orderID)]
	if !ok {
		return repository.ErrOrderNotFound
	}

	order.Status = status
	r.orders[uint64(orderID)] = order

	return nil
}
