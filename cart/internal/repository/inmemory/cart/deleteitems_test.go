package cart

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryDeleteItemsByUserID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		setupStore    func() map[int64]map[int64]uint16
		expectedStore map[int64]map[int64]uint16
		expectedCount int
	}{
		{
			name:   "clear existing user cart",
			userID: 1,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 5, 101: 3},
					2: {102: 2},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				2: {102: 2},
			},
			expectedCount: 2,
		},
		{
			name:   "clear non-existing user cart",
			userID: 3,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 5, 101: 3},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 5, 101: 3},
			},
			expectedCount: 0,
		},
		{
			name:   "clear empty cart store",
			userID: 1,
			setupStore: func() map[int64]map[int64]uint16 {
				return make(map[int64]map[int64]uint16)
			},
			expectedStore: map[int64]map[int64]uint16{},
			expectedCount: 0,
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

			count, err := repo.DeleteItemsByUserID(ctx, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStore, repo.store)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestRepositoryDeleteItemsByUserIDConcurrent(t *testing.T) {
	ctx := context.Background()
	repo := &Repository{
		mutex: sync.Mutex{},
		store: map[int64]map[int64]uint16{
			1: {100: 5, 101: 3},
			2: {102: 2},
		},
	}

	var wg sync.WaitGroup
	numRoutines := 100
	userID := int64(1)
	expectedCount := 2

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := repo.DeleteItemsByUserID(ctx, userID)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		}()
	}

	wg.Wait()
	expectedStore := map[int64]map[int64]uint16{
		2: {102: 2},
	}
	assert.Equal(t, expectedStore, repo.store)
}
