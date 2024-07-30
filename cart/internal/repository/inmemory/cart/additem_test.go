package cart

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCartAdd(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		skuID         int64
		count         uint16
		setupStore    func() map[int64]map[int64]uint16
		expectedStore map[int64]map[int64]uint16
	}{
		{
			name:   "add new item to new user",
			userID: 1,
			skuID:  100,
			count:  2,
			setupStore: func() map[int64]map[int64]uint16 {
				return make(map[int64]map[int64]uint16)
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 2},
			},
		},
		{
			name:   "add new item to existing user",
			userID: 1,
			skuID:  101,
			count:  3,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 2},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 2, 101: 3},
			},
		},
		{
			name:   "add existing item to user",
			userID: 1,
			skuID:  100,
			count:  3,
			setupStore: func() map[int64]map[int64]uint16 {
				return map[int64]map[int64]uint16{
					1: {100: 2},
				}
			},
			expectedStore: map[int64]map[int64]uint16{
				1: {100: 5},
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

			err := repo.Add(ctx, tt.userID, tt.skuID, tt.count)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStore, repo.store)
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	ctx := context.Background()

	repo := &Repository{
		store: make(map[int64]map[int64]uint16),
		mutex: sync.Mutex{},
	}

	userID := int64(1)
	skuID := int64(100)
	count := uint16(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.Add(ctx, userID, skuID, count)
	}
}

func TestRepositoryCartAddConcurrent(t *testing.T) {
	ctx := context.Background()
	repo := &Repository{
		mutex: sync.Mutex{},
		store: make(map[int64]map[int64]uint16),
	}

	var wg sync.WaitGroup
	numRoutines := 100
	userID := int64(1)
	skuID := int64(100)
	count := uint16(1)

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.Add(ctx, userID, skuID, count)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	expectedCount := uint16(numRoutines)
	assert.Equal(t, expectedCount, repo.store[userID][skuID])
}
