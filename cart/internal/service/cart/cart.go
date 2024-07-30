package cart

import (
	"context"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	lomsserviceModels "github.com/BruteMors/marketplace-service/cart/pkg/lomsservice/models"
	productserviceModels "github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
)

type CartRepository interface {
	Add(ctx context.Context, userID int64, skuID int64, count uint16) error
	DeleteItem(ctx context.Context, userID int64, skuID int64) error
	DeleteItemsByUserID(_ context.Context, userID int64) (int, error)
	GetCart(ctx context.Context, userID int64) ([]models.ItemCount, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, sku int64) (*productserviceModels.GetProductResponse, error)
	GetProducts(ctx context.Context, skus []int64) ([]productserviceModels.GetProductResponse, error)
	GetListSkus(ctx context.Context, startAfterSku, count int64) (
		*productserviceModels.ListSkusResponse,
		error,
	)
}

type LomsService interface {
	OrderCreate(ctx context.Context, order lomsserviceModels.OrderCreate) (orderID int64, err error)
	StocksInfo(ctx context.Context, sku uint32) (count uint64, err error)
}

type Service struct {
	productService ProductService
	cartRepository CartRepository
	lomsService    LomsService
}

func NewCartService(
	productService ProductService,
	repo CartRepository,
	lomsService LomsService,
) *Service {
	return &Service{
		productService: productService,
		cartRepository: repo,
		lomsService:    lomsService,
	}
}
