package cart

import (
	"context"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	AddItem(ctx context.Context, userID int64, skuID int64, count uint16) error
	DeleteItem(ctx context.Context, userID int64, skuID int64) error
	DeleteItemsByUserID(ctx context.Context, userID int64) error
	GetCart(ctx context.Context, userID int64) (*models.Cart, error)
	Checkout(ctx context.Context, userID int64) (orderID int64, err error)
}

type HttpApi struct {
	cartService Service
	validator   *validator.Validate
}

func NewCartHttpApi(
	cartService Service,
	validator *validator.Validate,
) *HttpApi {
	return &HttpApi{
		cartService: cartService,
		validator:   validator,
	}
}
