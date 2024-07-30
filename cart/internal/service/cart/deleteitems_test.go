package cart

import (
	"context"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/service/cart/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestServiceCartDeleteItemsByUserID(t *testing.T) {
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
		mockCartFunc  func()
		expectedCount int
		expectedError error
	}{
		{
			name:   "successful deletion",
			userID: 1,
			mockCartFunc: func() {
				cartRepositoryMock.DeleteItemsByUserIDMock.Expect(minimock.AnyContext, 1).Return(2, nil)
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:   "user cart not found",
			userID: 2,
			mockCartFunc: func() {
				cartRepositoryMock.DeleteItemsByUserIDMock.Expect(minimock.AnyContext, 2).Return(0, nil)
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:   "cart is empty",
			userID: 3,
			mockCartFunc: func() {
				cartRepositoryMock.DeleteItemsByUserIDMock.Expect(minimock.AnyContext, 3).Return(0, nil)
			},
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCartFunc()

			err := s.DeleteItemsByUserID(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
