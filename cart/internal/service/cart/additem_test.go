package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/service/cart/mock"
	productServiceModels "github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceCartAddItem(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	productServiceMock := mock.NewProductServiceMock(mc)
	cartRepositoryMock := mock.NewCartRepositoryMock(mc)
	lomsServiceMock := mock.NewLomsServiceMock(mc)
	s := &Service{
		productService: productServiceMock,
		cartRepository: cartRepositoryMock,
		lomsService:    lomsServiceMock,
	}

	ctx := context.Background()

	tests := []struct {
		name            string
		userID          int64
		skuID           int64
		count           uint16
		mockProductFunc func()
		mockLomsFunc    func()
		mockCartFunc    func()
		expectedError   error
	}{
		{
			name:   "success case",
			userID: 1,
			skuID:  100,
			count:  2,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 100).Return(
					&productServiceModels.GetProductResponse{Name: "Some product", Price: 100},
					nil,
				)
			},
			mockLomsFunc: func() {
				lomsServiceMock.StocksInfoMock.Expect(minimock.AnyContext, uint32(100)).Return(uint64(10), nil)
			},
			mockCartFunc: func() {
				cartRepositoryMock.AddMock.Expect(minimock.AnyContext, 1, 100, 2).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "product not found",
			userID: 1,
			skuID:  200,
			count:  2,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 200).Return(nil, productServiceModels.ErrNotFound)
			},
			mockLomsFunc:  func() {},
			mockCartFunc:  func() {},
			expectedError: models.ErrProductNotFound,
		},
		{
			name:   "product service error",
			userID: 1,
			skuID:  300,
			count:  2,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 300).Return(nil, errors.New("service error"))
			},
			mockLomsFunc:  func() {},
			mockCartFunc:  func() {},
			expectedError: errors.New("service error"),
		},
		{
			name:   "insufficient stocks",
			userID: 1,
			skuID:  100,
			count:  10,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 100).Return(
					&productServiceModels.GetProductResponse{Name: "Product", Price: 200},
					nil,
				)
			},
			mockLomsFunc: func() {
				lomsServiceMock.StocksInfoMock.Expect(minimock.AnyContext, uint32(100)).Return(uint64(5), nil)
			},
			mockCartFunc:  func() {},
			expectedError: models.ErrStocksNotEnough,
		},
		{
			name:   "loms service error",
			userID: 1,
			skuID:  200,
			count:  2,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 200).Return(
					&productServiceModels.GetProductResponse{Name: "Product", Price: 100},
					nil,
				)
			},
			mockLomsFunc: func() {
				lomsServiceMock.StocksInfoMock.Expect(minimock.AnyContext, uint32(200)).Return(uint64(0), errors.New("loms service error"))
			},
			mockCartFunc:  func() {},
			expectedError: errors.New("loms service error"),
		},
		{
			name:   "cart repository add error",
			userID: 1,
			skuID:  400,
			count:  2,
			mockProductFunc: func() {
				productServiceMock.GetProductMock.Expect(minimock.AnyContext, 400).Return(
					&productServiceModels.GetProductResponse{Name: "Some product", Price: 100},
					nil,
				)
			},
			mockLomsFunc: func() {
				lomsServiceMock.StocksInfoMock.Expect(minimock.AnyContext, uint32(400)).Return(uint64(3), nil)
			},
			mockCartFunc: func() {
				cartRepositoryMock.AddMock.Expect(minimock.AnyContext, 1, 400, 2).Return(errors.New("add error"))
			},
			expectedError: errors.New("add error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockProductFunc()
			tt.mockLomsFunc()
			tt.mockCartFunc()

			err := s.AddItem(ctx, tt.userID, tt.skuID, tt.count)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
