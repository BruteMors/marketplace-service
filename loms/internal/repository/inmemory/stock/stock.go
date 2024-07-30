package stock

import (
	_ "embed"
	"encoding/json"
	"sync"

	stockdomain "github.com/BruteMors/marketplace-service/loms/internal/domain/stock"
)

//go:embed stock-data.json
var stockData []byte

type Repository struct {
	mu    sync.RWMutex
	stock map[uint32]stockdomain.Item
}

func NewRepository() (*Repository, error) {
	repo := &Repository{}

	var tempItems []struct {
		SKU        uint32 `json:"sku"`
		TotalCount uint64 `json:"total_count"`
		Reserved   uint64 `json:"reserved"`
	}

	if err := json.Unmarshal(stockData, &tempItems); err != nil {
		return nil, err
	}

	repo.stock = make(map[uint32]stockdomain.Item, len(tempItems))

	for _, item := range tempItems {
		repo.stock[item.SKU] = stockdomain.Item{
			SKU:        item.SKU,
			TotalCount: item.TotalCount,
			Reserved:   item.Reserved,
		}
	}

	return repo, nil
}
