package tests

import (
	"context"
	"testing"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/stretchr/testify/require"
)

func TestSetOrderStatus(t *testing.T) {
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

	newStatus := ordermodels.Status("cancelled")
	err = repo.SetStatus(ctx, orderID, newStatus)
	require.NoError(t, err)

	updatedOrder, err := repo.GetByID(ctx, orderID)
	require.NoError(t, err)

	require.Equal(t, newStatus, updatedOrder.Status)
}
