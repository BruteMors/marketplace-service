package order

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/models"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/responses"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/service/order/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceOrderInfo(t *testing.T) {
	mc := minimock.NewController(t)
	orderRepositoryMock := mock.NewRepositoryMock(mc)
	s := &Service{
		orderRepository: orderRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name             string
		orderID          int64
		mockOrderFunc    func()
		expectedResponse responses.OrderInfo
		expectedError    error
	}{
		{
			name:    "order not found",
			orderID: 1,
			mockOrderFunc: func() {
				orderRepositoryMock.GetByIDMock.Expect(minimock.AnyContext, 1).Return(order.Order{}, repository.ErrOrderNotFound)
			},
			expectedResponse: responses.OrderInfo{},
			expectedError:    models.ErrOrderNotFound,
		},
		{
			name:    "database error",
			orderID: 2,
			mockOrderFunc: func() {
				orderRepositoryMock.GetByIDMock.Expect(minimock.AnyContext, 2).Return(order.Order{}, errors.New("database error"))
			},
			expectedResponse: responses.OrderInfo{},
			expectedError:    errors.New("database error"),
		},
		{
			name:    "successful info retrieval",
			orderID: 3,
			mockOrderFunc: func() {
				order := &order.Order{
					Status: order.OrderStatusNew,
					UserID: 123,
					Items: []order.Item{
						{SKU: 100, Count: 2},
						{SKU: 101, Count: 1},
					},
				}
				orderRepositoryMock.GetByIDMock.Expect(minimock.AnyContext, 3).Return(*order, nil)
			},
			expectedResponse: responses.OrderInfo{
				Status: order.OrderStatusNew,
				User:   123,
				Items: []responses.Item{
					{SKU: 100, Count: 2},
					{SKU: 101, Count: 1},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockOrderFunc()
			response, err := s.OrderInfo(ctx, tt.orderID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedResponse, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}
