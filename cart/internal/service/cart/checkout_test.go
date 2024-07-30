package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	"github.com/BruteMors/marketplace-service/cart/internal/service/cart/mock"
	lomsserviceModels "github.com/BruteMors/marketplace-service/cart/pkg/lomsservice/models"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceCheckout(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	cartRepositoryMock := mock.NewCartRepositoryMock(mc)
	lomsServiceMock := mock.NewLomsServiceMock(mc)
	s := &Service{
		cartRepository: cartRepositoryMock,
		lomsService:    lomsServiceMock,
	}

	ctx := context.Background()

	tests := []struct {
		name               string
		userID             int64
		mockCartFunc       func()
		mockLomsFunc       func()
		mockDeleteCartFunc func()
		expectedOrderID    int64
		expectedError      error
	}{
		{
			name:   "cart not found",
			userID: 1,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 1).Return(nil, repository.ErrCartNotFound)
			},
			expectedOrderID: 0,
			expectedError:   models.ErrCartNotFound,
		},
		{
			name:   "cart repository error",
			userID: 2,
			mockCartFunc: func() {
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 2).Return(nil, errors.New("db error"))
			},
			expectedOrderID: 0,
			expectedError:   errors.New("db error"),
		},
		{
			name:   "successful checkout",
			userID: 3,
			mockCartFunc: func() {
				cart := []models.ItemCount{
					{SkuID: 100, Count: 2},
					{SkuID: 101, Count: 3},
				}
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 3).Return(cart, nil)
			},
			mockLomsFunc: func() {
				items := []lomsserviceModels.OrderItem{
					{SkuID: 100, Count: 2},
					{SkuID: 101, Count: 3},
				}
				lomsServiceMock.OrderCreateMock.Expect(minimock.AnyContext, lomsserviceModels.OrderCreate{
					User:  3,
					Items: items,
				}).Return(int64(12345), nil)
			},
			mockDeleteCartFunc: func() {
				cartRepositoryMock.DeleteItemsByUserIDMock.Expect(minimock.AnyContext, 3).Return(2, nil)
			},
			expectedOrderID: 12345,
			expectedError:   nil,
		},
		{
			name:   "loms service error",
			userID: 4,
			mockCartFunc: func() {
				cart := []models.ItemCount{
					{SkuID: 200, Count: 1},
				}
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 4).Return(cart, nil)
			},
			mockLomsFunc: func() {
				items := []lomsserviceModels.OrderItem{
					{SkuID: 200, Count: 1},
				}
				lomsServiceMock.OrderCreateMock.Expect(minimock.AnyContext, lomsserviceModels.OrderCreate{
					User:  4,
					Items: items,
				}).Return(int64(0), errors.New("loms service error"))
			},
			expectedOrderID: 0,
			expectedError:   errors.New("loms service error"),
		},
		{
			name:   "delete cart items error",
			userID: 5,
			mockCartFunc: func() {
				cart := []models.ItemCount{
					{SkuID: 300, Count: 4},
				}
				cartRepositoryMock.GetCartMock.Expect(minimock.AnyContext, 5).Return(cart, nil)
			},
			mockLomsFunc: func() {
				items := []lomsserviceModels.OrderItem{
					{SkuID: 300, Count: 4},
				}
				lomsServiceMock.OrderCreateMock.Expect(minimock.AnyContext, lomsserviceModels.OrderCreate{
					User:  5,
					Items: items,
				}).Return(int64(67890), nil)
			},
			mockDeleteCartFunc: func() {
				cartRepositoryMock.DeleteItemsByUserIDMock.Expect(minimock.AnyContext, 5).Return(0, errors.New("delete error"))
			},
			expectedOrderID: 0,
			expectedError:   errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCartFunc()
			if tt.mockLomsFunc != nil {
				tt.mockLomsFunc()
			}
			if tt.mockDeleteCartFunc != nil {
				tt.mockDeleteCartFunc()
			}

			orderID, err := s.Checkout(ctx, tt.userID)
			assert.Equal(t, tt.expectedOrderID, orderID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
