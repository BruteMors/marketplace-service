package responses

import "github.com/BruteMors/marketplace-service/loms/internal/models/order"

type OrderInfo struct {
	Status order.Status
	User   int64
	Items  []Item
}

type Item struct {
	SKU   uint32
	Count uint16
}
