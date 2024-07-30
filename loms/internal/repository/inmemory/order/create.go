package order

import (
	"context"
	"time"

	orderdomain "github.com/BruteMors/marketplace-service/loms/internal/domain/order"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
)

func (r *Repository) Create(ctx context.Context, order ordermodels.NewOrder) (orderID int64, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderID = int64(r.ordersCounter)
	r.ordersCounter++

	items := make([]orderdomain.Item, 0, len(order.Items))

	for _, i := range order.Items {
		items = append(items, orderdomain.Item{
			SKU:   i.SKU,
			Count: i.Count,
		})
	}

	r.orders[uint64(orderID)] = orderdomain.Order{
		ID:        orderID,
		Status:    order.Status,
		UserID:    order.User,
		Items:     items,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
	}

	return orderID, nil
}
