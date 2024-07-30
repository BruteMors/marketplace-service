package order

import (
	"time"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
)

type Order struct {
	ID        int64
	Status    ordermodels.Status
	UserID    int64
	Items     []Item
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Item struct {
	SKU   uint32
	Count uint16
}
