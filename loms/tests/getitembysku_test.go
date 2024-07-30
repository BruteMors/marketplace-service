package tests

import (
	"context"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/stretchr/testify/require"
)

func TestGetItemBySKU(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()

	client, err := pg.New(ctx, pgMasterDSN, pgReplicaDSNs)
	require.NoError(t, err)

	repo := stock.NewRepository(client)

	skuID := uint32(1076963)

	item, err := repo.GetBySKU(ctx, skuID)
	require.NoError(t, err)

	require.Equal(t, skuID, item.SKU)
	require.Equal(t, uint64(100), item.TotalCount)
	require.Equal(t, uint64(0), item.Reserved)
}
