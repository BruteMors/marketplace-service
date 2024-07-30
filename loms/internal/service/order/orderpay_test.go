package order

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/models"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/service/order/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceOrderPay(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepositoryMock := mock.NewRepositoryMock(mc)
	stockServiceMock := mock.NewStockServiceMock(mc)
	statusOutboxRepositoryMock := mock.NewStatusOutboxRepositoryMock(mc)
	txManagerMock := mock.NewTxManagerMock(mc)

	s := &Service{
		orderRepository:        orderRepositoryMock,
		stockService:           stockServiceMock,
		statusOutboxRepository: statusOutboxRepositoryMock,
		txManager:              txManagerMock,
	}

	ctx := context.Background()

	tests := []struct {
		name              string
		orderID           int64
		mockOrderFunc     func()
		mockReserveFunc   func()
		mockSetStatusFunc func()
		mockStatusOutbox  func()
		expectedError     error
	}{
		{
			name:    "order not found",
			orderID: 1,
			mockOrderFunc: func() {
				orderRepositoryMock.GetByIDMock.Expect(ctx, 1).Return(ordermodels.Order{}, repository.ErrOrderNotFound)
			},
			expectedError: models.ErrOrderNotFound,
		},
		{
			name:    "order repository error",
			orderID: 2,
			mockOrderFunc: func() {
				orderRepositoryMock.GetByIDMock.Expect(ctx, 2).Return(ordermodels.Order{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:    "successful payment",
			orderID: 3,
			mockOrderFunc: func() {
				order := ordermodels.Order{
					Items: []ordermodels.Item{
						{SKU: 100, Count: 2},
						{SKU: 101, Count: 1},
					},
				}
				orderRepositoryMock.GetByIDMock.Expect(ctx, 3).Return(order, nil)
			},
			mockReserveFunc: func() {
				items := []stock.ReserveItem{
					{SKU: 100, Count: 2},
					{SKU: 101, Count: 1},
				}
				stockServiceMock.ReserveRemoveMock.Expect(ctx, items).Return(nil)
			},
			mockSetStatusFunc: func() {
				orderRepositoryMock.SetStatusMock.Expect(ctx, 3, ordermodels.OrderStatusPayed).Return(nil)
			},
			mockStatusOutbox: func() {
				statusOutboxRepositoryMock.CreateOrderStatusChangedEventMock.Expect(ctx, 3, ordermodels.OrderStatusPayed).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "reserve removal error",
			orderID: 4,
			mockOrderFunc: func() {
				order := ordermodels.Order{
					Items: []ordermodels.Item{
						{SKU: 200, Count: 1},
					},
				}
				orderRepositoryMock.GetByIDMock.Expect(ctx, 4).Return(order, nil)
			},
			mockReserveFunc: func() {
				items := []stock.ReserveItem{
					{SKU: 200, Count: 1},
				}
				stockServiceMock.ReserveRemoveMock.Expect(ctx, items).Return(errors.New("reserve removal error"))
			},
			expectedError: errors.New("reserve removal error"),
		},
		{
			name:    "order status update error",
			orderID: 5,
			mockOrderFunc: func() {
				order := ordermodels.Order{
					Items: []ordermodels.Item{
						{SKU: 300, Count: 4},
					},
				}
				orderRepositoryMock.GetByIDMock.Expect(ctx, 5).Return(order, nil)
			},
			mockReserveFunc: func() {
				items := []stock.ReserveItem{
					{SKU: 300, Count: 4},
				}
				stockServiceMock.ReserveRemoveMock.Expect(ctx, items).Return(nil)
			},
			mockSetStatusFunc: func() {
				orderRepositoryMock.SetStatusMock.Expect(ctx, 5, ordermodels.OrderStatusPayed).Return(errors.New("status update error"))
			},
			expectedError: errors.New("status update error"),
		},
		{
			name:    "status outbox error",
			orderID: 6,
			mockOrderFunc: func() {
				order := ordermodels.Order{
					Items: []ordermodels.Item{
						{SKU: 400, Count: 1},
					},
				}
				orderRepositoryMock.GetByIDMock.Expect(ctx, 6).Return(order, nil)
			},
			mockReserveFunc: func() {
				items := []stock.ReserveItem{
					{SKU: 400, Count: 1},
				}
				stockServiceMock.ReserveRemoveMock.Expect(ctx, items).Return(nil)
			},
			mockSetStatusFunc: func() {
				orderRepositoryMock.SetStatusMock.Expect(ctx, 6, ordermodels.OrderStatusPayed).Return(nil)
			},
			mockStatusOutbox: func() {
				statusOutboxRepositoryMock.CreateOrderStatusChangedEventMock.Expect(ctx, 6, ordermodels.OrderStatusPayed).Return(errors.New("status outbox error"))
			},
			expectedError: errors.New("status outbox error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockOrderFunc()
			if tt.mockReserveFunc != nil {
				tt.mockReserveFunc()
			}
			if tt.mockSetStatusFunc != nil {
				tt.mockSetStatusFunc()
			}
			if tt.mockStatusOutbox != nil {
				tt.mockStatusOutbox()
			}

			err := s.orderPay(ctx, tt.orderID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
