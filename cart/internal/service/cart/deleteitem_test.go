package cart

import (
	"context"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/service/cart/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceCartDeleteItem(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	cartRepositoryMock := mock.NewCartRepositoryMock(mc)
	s := &Service{
		cartRepository: cartRepositoryMock,
	}

	ctx := context.Background()

	tests := []struct {
		name          string
		userID        int64
		skuID         int64
		mockCartFunc  func()
		expectedError error
	}{
		{
			name:   "successful deletion",
			userID: 1,
			skuID:  100,
			mockCartFunc: func() {
				cartRepositoryMock.DeleteItemMock.Expect(minimock.AnyContext, 1, 100).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "item not found",
			userID: 1,
			skuID:  200,
			mockCartFunc: func() {
				cartRepositoryMock.DeleteItemMock.Expect(minimock.AnyContext, 1, 200).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCartFunc()

			err := s.DeleteItem(ctx, tt.userID, tt.skuID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
