package stock

import (
	"context"
	"errors"
	"testing"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/service/stock/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceReserve(t *testing.T) {
	mc := minimock.NewController(t)
	stockRepositoryMock := mock.NewRepositoryMock(mc)
	s := &Service{
		stockRepository: stockRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name            string
		items           []ordermodels.Item
		mockReserveFunc func()
		expectedError   error
	}{
		{
			name: "successful reserve",
			items: []ordermodels.Item{
				{SKU: 100, Count: 10},
				{SKU: 101, Count: 5},
			},
			mockReserveFunc: func() {
				reserveItems := []stock.ReserveItem{
					{SKU: 100, Count: 10},
					{SKU: 101, Count: 5},
				}
				stockRepositoryMock.ReserveMock.Expect(minimock.AnyContext, reserveItems).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "reserve fails due to database error",
			items: []ordermodels.Item{
				{SKU: 200, Count: 3},
			},
			mockReserveFunc: func() {
				reserveItems := []stock.ReserveItem{
					{SKU: 200, Count: 3},
				}
				stockRepositoryMock.ReserveMock.Expect(minimock.AnyContext, reserveItems).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:  "empty items list",
			items: []ordermodels.Item{},
			mockReserveFunc: func() {
				reserveItems := []stock.ReserveItem{}
				stockRepositoryMock.ReserveMock.Expect(minimock.AnyContext, reserveItems).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReserveFunc()
			err := s.Reserve(ctx, tt.items)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
