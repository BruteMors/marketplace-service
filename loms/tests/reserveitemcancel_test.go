package tests

import (
	"context"
	"testing"

	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/stretchr/testify/require"
)

func TestReserveCancel(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()

	client, err := pg.New(ctx, pgMasterDSN, pgReplicaDSNs)
	require.NoError(t, err)

	repo := stock.NewRepository(client)

	items := []stockmodels.ReserveItem{
		{SKU: 1076963, Count: 10},
		{SKU: 1148162, Count: 20},
	}

	err = repo.Reserve(ctx, items)
	require.NoError(t, err)

	var reservedCount int
	err = client.MasterDB().QueryRow(ctx, "SELECT reserved FROM items WHERE sku=$1", items[0].SKU).Scan(&reservedCount)
	require.NoError(t, err)
	require.Equal(t, 10, reservedCount)

	err = client.MasterDB().QueryRow(ctx, "SELECT reserved FROM items WHERE sku=$1", items[1].SKU).Scan(&reservedCount)
	require.NoError(t, err)
	require.Equal(t, 20, reservedCount)

	err = repo.ReserveCancel(ctx, items)
	require.NoError(t, err)

	err = client.MasterDB().QueryRow(ctx, "SELECT reserved FROM items WHERE sku=$1", items[0].SKU).Scan(&reservedCount)
	require.NoError(t, err)
	require.Equal(t, 0, reservedCount)

	err = client.MasterDB().QueryRow(ctx, "SELECT reserved FROM items WHERE sku=$1", items[1].SKU).Scan(&reservedCount)
	require.NoError(t, err)
	require.Equal(t, 0, reservedCount)
}
