package order

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/models"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/requests"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/service/order/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceOrderCreate(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepositoryMock := mock.NewRepositoryMock(mc)
	stockServiceMock := mock.NewStockServiceMock(mc)
	s := &Service{
		orderRepository: orderRepositoryMock,
		stockService:    stockServiceMock,
	}

	ctx := context.Background()

	tests := []struct {
		name                string
		request             *requests.OrderCreate
		mockOrderCreateFunc func()
		mockReserveFunc     func()
		mockSetStatusFunc   func()
		expectedOrderID     int64
		expectedError       error
	}{
		{
			name: "successful order creation",
			request: &requests.OrderCreate{
				User: 1,
				Items: []requests.Item{
					{SKU: 100, Count: 2},
					{SKU: 101, Count: 3},
				},
			},
			mockOrderCreateFunc: func() {
				items := []ordermodels.Item{
					{SKU: 100, Count: 2},
					{SKU: 101, Count: 3},
				}
				newOrder := ordermodels.NewOrder{
					User:   1,
					Items:  items,
					Status: ordermodels.OrderStatusNew,
				}
				orderRepositoryMock.CreateMock.Expect(ctx, newOrder).Return(int64(12345), nil)
			},
			mockReserveFunc: func() {
				items := []ordermodels.Item{
					{SKU: 100, Count: 2},
					{SKU: 101, Count: 3},
				}
				stockServiceMock.ReserveMock.Expect(ctx, items).Return(nil)
			},
			mockSetStatusFunc: func() {
				orderRepositoryMock.SetStatusMock.Expect(ctx, 12345, ordermodels.OrderStatusAwaitingPayment).Return(nil)
			},
			expectedOrderID: 12345,
			expectedError:   nil,
		},
		{
			name: "order creation fails",
			request: &requests.OrderCreate{
				User: 2,
				Items: []requests.Item{
					{SKU: 200, Count: 1},
				},
			},
			mockOrderCreateFunc: func() {
				items := []ordermodels.Item{
					{SKU: 200, Count: 1},
				}
				newOrder := ordermodels.NewOrder{
					User:   2,
					Items:  items,
					Status: ordermodels.OrderStatusNew,
				}
				orderRepositoryMock.CreateMock.Expect(ctx, newOrder).Return(int64(0), errors.New("db error"))
			},
			expectedOrderID: 0,
			expectedError:   errors.New("db error"),
		},
		{
			name: "reserve fails with SKU not found",
			request: &requests.OrderCreate{
				User: 3,
				Items: []requests.Item{
					{SKU: 300, Count: 4},
				},
			},
			mockOrderCreateFunc: func() {
				items := []ordermodels.Item{
					{SKU: 300, Count: 4},
				}
				newOrder := ordermodels.NewOrder{
					User:   3,
					Items:  items,
					Status: ordermodels.OrderStatusNew,
				}
				orderRepositoryMock.CreateMock.Expect(ctx, newOrder).Return(int64(67890), nil)
			},
			mockReserveFunc: func() {
				items := []ordermodels.Item{
					{SKU: 300, Count: 4},
				}
				stockServiceMock.ReserveMock.Expect(ctx, items).Return(repository.ErrSKUNotFound)
			},
			mockSetStatusFunc: func() {
				orderRepositoryMock.SetStatusMock.Expect(ctx, 67890, ordermodels.OrderStatusFailed).Return(nil)
			},
			expectedOrderID: 0,
			expectedError:   models.ErrSKUNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockOrderCreateFunc()
			if tt.mockReserveFunc != nil {
				tt.mockReserveFunc()
			}
			if tt.mockSetStatusFunc != nil {
				tt.mockSetStatusFunc()
			}

			orderID, err := s.orderCreate(ctx, tt.request)
			assert.Equal(t, tt.expectedOrderID, orderID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
