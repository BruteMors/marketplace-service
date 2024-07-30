package cart

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryDeleteItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		skuID         int64
		setupStore    func() map[int64]map[int64]uint16
		expectedStore map[int64]map[int64]uint16
	}{
		{
			name:   "delete existing item",
			userID: 1,
			skuID:  100,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 5, 101: 3},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {101: 3},
			},
		},
		{
			name:   "delete non-existing item from user with other items",
			userID: 1,
			skuID:  102,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 5, 101: 3},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 5, 101: 3},
			},
		},
		{
			name:   "delete item from non-existing user",
			userID: 2,
			skuID:  100,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 5, 101: 3},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 5, 101: 3},
			},
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

			err := repo.DeleteItem(ctx, tt.userID, tt.skuID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStore, repo.store)
		})
	}
}

func TestRepositoryDeleteItemConcurrent(t *testing.T) {
	ctx := context.Background()
	repo := &Repository{
		mutex: sync.Mutex{},
		store: map[int64]map[int64]uint16{
			1: {100: 100},
		},
	}

	var wg sync.WaitGroup
	numRoutines := 100
	userID := int64(1)
	skuID := int64(100)

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.DeleteItem(ctx, userID, skuID)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	expectedStore := map[int64]map[int64]uint16{
		1: {},
	}
	assert.Equal(t, expectedStore, repo.store)
}
