package stock

import (
	"context"
	"errors"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/models"
	"github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/service/stock/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceStocksInfo(t *testing.T) {
	mc := minimock.NewController(t)
	stockRepositoryMock := mock.NewRepositoryMock(mc)
	s := &Service{
		stockRepository: stockRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name           string
		skuID          uint32
		mockStocksFunc func()
		expectedCount  uint64
		expectedError  error
	}{
		{
			name:  "successful stocks info retrieval",
			skuID: 100,
			mockStocksFunc: func() {
				item := stock.Item{
					SKU:        100,
					TotalCount: 15,
					Reserved:   5,
				}
				stockRepositoryMock.GetBySKUMock.Expect(minimock.AnyContext, 100).Return(item, nil)
			},
			expectedCount: 10,
			expectedError: nil,
		},
		{
			name:  "SKU not found",
			skuID: 101,
			mockStocksFunc: func() {
				stockRepositoryMock.GetBySKUMock.Expect(minimock.AnyContext, 101).Return(stock.Item{}, repository.ErrSKUNotFound)
			},
			expectedCount: 0,
			expectedError: models.ErrSKUNotFound,
		},
		{
			name:  "database error",
			skuID: 102,
			mockStocksFunc: func() {
				stockRepositoryMock.GetBySKUMock.Expect(minimock.AnyContext, 102).Return(stock.Item{}, errors.New("database error"))
			},
			expectedCount: 0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockStocksFunc()
			count, err := s.StocksInfo(ctx, tt.skuID)
			assert.Equal(t, tt.expectedCount, count)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
