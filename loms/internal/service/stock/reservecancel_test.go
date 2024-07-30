package stock

import (
	"context"
	"errors"
	"testing"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/service/stock/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceReserveCancel(t *testing.T) {
	mc := minimock.NewController(t)
	stockRepositoryMock := mock.NewRepositoryMock(mc)
	s := &Service{
		stockRepository: stockRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name           string
		items          []stockmodels.ReserveItem
		mockCancelFunc func()
		expectedError  error
	}{
		{
			name: "successful reserve cancel",
			items: []stockmodels.ReserveItem{
				{SKU: 100, Count: 10},
				{SKU: 101, Count: 5},
			},
			mockCancelFunc: func() {
				stockRepositoryMock.ReserveCancelMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{
					{SKU: 100, Count: 10},
					{SKU: 101, Count: 5},
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "reserve cancel fails due to database error",
			items: []stockmodels.ReserveItem{
				{SKU: 200, Count: 3},
			},
			mockCancelFunc: func() {
				stockRepositoryMock.ReserveCancelMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{
					{SKU: 200, Count: 3},
				}).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:  "empty items list",
			items: []stockmodels.ReserveItem{},
			mockCancelFunc: func() {
				stockRepositoryMock.ReserveCancelMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCancelFunc()
			err := s.ReserveCancel(ctx, tt.items)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
