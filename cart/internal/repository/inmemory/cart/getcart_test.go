package cart

import (
	"context"
	"sync"
	"testing"

	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryGetCart(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		setupStore    func() map[int64]map[int64]uint16
		expectedItems []models.ItemCount
		expectedError error
	}{
		{
			name:   "existing user with non-empty cart",
			userID: 1,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {101: 3, 100: 2},
				}
			},
			expectedItems: []models.ItemCount{
				{SkuID: 101, Count: 3},
				{SkuID: 100, Count: 2},
			},
			expectedError: nil,
		},
		{
			name:   "existing user with empty cart",
			userID: 1,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {},
				}
			},
			expectedItems: nil,
			expectedError: repository.ErrCartEmpty,
		},
		{
			name:   "non-existing user",
			userID: 2,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 2, 101: 3},
				}
			},
			expectedItems: nil,
			expectedError: repository.ErrCartNotFound,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &Repository{
				mutex: sync.Mutex{},
				store: tt.setupStore(),
			}

			items, err := repo.GetCart(ctx, tt.userID)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedItems, items)
			}
		})
	}
}

func TestRepositoryGetCartConcurrent(t *testing.T) {
	ctx := context.Background()
	repo := &Repository{
		mutex: sync.Mutex{},
		store: map[int64]map[int64]uint16{
			1: {101: 3, 100: 2},
			2: {102: 5},
		},
	}

	var wg sync.WaitGroup
	numRoutines := 100
	userID := int64(1)

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			items, err := repo.GetCart(ctx, userID)
			assert.NoError(t, err)
			assert.ElementsMatch(t, []models.ItemCount{
				{SkuID: 101, Count: 3},
				{SkuID: 100, Count: 2},
			}, items)
		}()
	}

	wg.Wait()
}
