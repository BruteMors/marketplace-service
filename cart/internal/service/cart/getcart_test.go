package cart

import (
	"context"
	"fmt"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	"github.com/BruteMors/marketplace-service/cart/internal/service/cart/mock"
	productServiceModels "github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceCartGetCart(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	cartRepositoryMock := mock.NewCartRepositoryMock(mc)
	productServiceMock := mock.NewProductServiceMock(mc)
	s := &Service{
		cartRepository: cartRepositoryMock,
		productService: productServiceMock,
	}

	ctx := context.Background()

	tests := []struct {
		name            string
		userID          int64
		mockCartFunc    func()
		mockProductFunc func()
		expectedResult  *models.Cart
		expectedError   error
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 1).Return([]models.ItemCount{
					{SkuID: 101, Count: 2},
					{SkuID: 102, Count: 3},
				}, nil)
			},
			mockProductFunc: func() {
				productServiceMock.GetProductsMock.Expect(minimock.AnyContext, []int64{101, 102}).Return(
					[]productServiceModels.GetProductResponse{
						{Sku: 101, Name: "Product A", Price: 100},
						{Sku: 102, Name: "Product B", Price: 150},
					},
					nil,
				)
			},
			expectedResult: &models.Cart{
				Items: []models.Item{
					{Name: "Product A", Price: 100, ItemCount: models.ItemCount{SkuID: 101, Count: 2}},
					{Name: "Product B", Price: 150, ItemCount: models.ItemCount{SkuID: 102, Count: 3}},
				},
				TotalPrice: 650,
			},
			expectedError: nil,
		},
		{
			name:   "cart not found error",
			userID: 2,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 2).Return(nil, repository.ErrCartNotFound)
			},
			mockProductFunc: func() {},
			expectedResult:  nil,
			expectedError:   models.ErrCartNotFound,
		},
		{
			name:   "product not found error",
			userID: 3,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 3).Return(
					[]models.ItemCount{{SkuID: 103, Count: 1}},
					nil,
				)
			},
			mockProductFunc: func() {
				productServiceMock.GetProductsMock.Expect(minimock.AnyContext, []int64{103}).Return(
					nil,
					productServiceModels.ErrNotFound,
				)
			},
			expectedResult: nil,
			expectedError:  models.ErrProductNotFound,
		},
		{
			name:   "unexpected repository error",
			userID: 2,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 2).Return(
					nil,
					fmt.Errorf("unexpected repository error"),
				)
			},
			mockProductFunc: func() {},
			expectedResult:  nil,
			expectedError:   fmt.Errorf("unexpected repository error"),
		},
		{
			name:   "unexpected product service error",
			userID: 3,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 3).Return([]models.ItemCount{{SkuID: 103, Count: 1}}, nil)
			},
			mockProductFunc: func() {
				productServiceMock.GetProductsMock.Expect(minimock.AnyContext, []int64{103}).Return(
					nil,
					fmt.Errorf("unexpected product service error"),
				)
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("unexpected product service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCartFunc()
			tt.mockProductFunc()

			result, err := s.GetCart(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestServiceCartCalculateTotalPrice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		items         []models.Item
		expectedPrice uint32
	}{
		{
			name:          "empty cart",
			items:         []models.Item{},
			expectedPrice: 0,
		},
		{
			name: "single item",
			items: []models.Item{
				{Name: "Product A", Price: 100, ItemCount: models.ItemCount{SkuID: 101, Count: 1}},
			},
			expectedPrice: 100,
		},
		{
			name: "multiple items",
			items: []models.Item{
				{Name: "Product A", Price: 100, ItemCount: models.ItemCount{SkuID: 101, Count: 2}},
				{Name: "Product B", Price: 200, ItemCount: models.ItemCount{SkuID: 102, Count: 3}},
			},
			expectedPrice: 800,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			price := s.calculateTotalPrice(tt.items)
			assert.Equal(t, tt.expectedPrice, price)
		})
	}
}
