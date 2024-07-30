package order

import (
	"context"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (order ordermodels.Order, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, ok := r.orders[uint64(orderID)]
	if !ok {
		return ordermodels.Order{}, repository.ErrOrderNotFound
	}

	var items []ordermodels.Item
	for _, item := range repoOrder.Items {
		items = append(items, ordermodels.Item{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	order = ordermodels.Order{
		ID:        repoOrder.ID,
		Status:    repoOrder.Status,
		UserID:    repoOrder.UserID,
		Items:     items,
		CreatedAt: repoOrder.CreatedAt,
		UpdatedAt: repoOrder.UpdatedAt,
	}

	return order, nil
}
