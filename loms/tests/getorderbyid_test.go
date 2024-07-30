package tests

import (
	"context"
	"testing"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/stretchr/testify/require"
)

func TestGetOrderByID(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()

	client, err := pg.New(ctx, pgMasterDSN, pgReplicaDSNs)
	require.NoError(t, err)

	repo := order.NewRepository(client)

	newOrder := ordermodels.NewOrder{
		User:   1,
		Status: "new",
		Items: []ordermodels.Item{
			{SKU: 1076963, Count: 1},
			{SKU: 1148162, Count: 2},
		},
	}

	orderID, err := repo.Create(ctx, newOrder)
	require.NoError(t, err)
	require.NotZero(t, orderID)

	createdOrder, err := repo.GetByID(ctx, orderID)
	require.NoError(t, err)

	require.Equal(t, orderID, createdOrder.ID)
	require.Equal(t, int64(1), createdOrder.UserID)
	require.Equal(t, "new", string(createdOrder.Status))
	require.Len(t, createdOrder.Items, 2)

	expectedItems := []ordermodels.Item{
		{SKU: 1076963, Count: 1},
		{SKU: 1148162, Count: 2},
	}

	for i, item := range createdOrder.Items {
		require.Equal(t, expectedItems[i].SKU, item.SKU)
		require.Equal(t, expectedItems[i].Count, item.Count)
	}
}
