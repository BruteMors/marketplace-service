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

func TestServiceReserveRemove(t *testing.T) {
	mc := minimock.NewController(t)
	stockRepositoryMock := mock.NewRepositoryMock(mc)
	s := &Service{
		stockRepository: stockRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name           string
		items          []stockmodels.ReserveItem
		mockRemoveFunc func()
		expectedError  error
	}{
		{
			name: "successful reserve removal",
			items: []stockmodels.ReserveItem{
				{SKU: 100, Count: 10},
				{SKU: 101, Count: 5},
			},
			mockRemoveFunc: func() {
				stockRepositoryMock.ReserveRemoveMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{
					{SKU: 100, Count: 10},
					{SKU: 101, Count: 5},
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "reserve removal fails due to database error",
			items: []stockmodels.ReserveItem{
				{SKU: 200, Count: 3},
			},
			mockRemoveFunc: func() {
				stockRepositoryMock.ReserveRemoveMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{
					{SKU: 200, Count: 3},
				}).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name:  "empty items list",
			items: []stockmodels.ReserveItem{},
			mockRemoveFunc: func() {
				stockRepositoryMock.ReserveRemoveMock.Expect(minimock.AnyContext, []stockmodels.ReserveItem{}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRemoveFunc()
			err := s.ReserveRemove(ctx, tt.items)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
