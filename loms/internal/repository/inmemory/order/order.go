package order

import (
	_ "embed"
	"sync"

	orderdomain "github.com/BruteMors/marketplace-service/loms/internal/domain/order"
)

type Repository struct {
	mu            sync.RWMutex
	ordersCounter uint64
	orders        map[uint64]orderdomain.Order
}

func NewRepository() (*Repository, error) {
	repo := &Repository{
		orders: make(map[uint64]orderdomain.Order),
	}

	return repo, nil
}
