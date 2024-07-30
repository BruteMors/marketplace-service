package tests

import (
	"context"
	"testing"

	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
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

	var createdOrder ordermodels.Order
	err = client.MasterDB().QueryRow(ctx, "SELECT order_id, user_id, status FROM orders WHERE order_id=$1", orderID).Scan(&createdOrder.ID, &createdOrder.UserID, &createdOrder.Status)
	require.NoError(t, err)
	require.Equal(t, int64(1), createdOrder.UserID)
	require.Equal(t, "new", string(createdOrder.Status))

	var itemCount int
	err = client.MasterDB().QueryRow(ctx, "SELECT COUNT(*) FROM orders_to_items WHERE order_id=$1", orderID).Scan(&itemCount)
	require.NoError(t, err)
	require.Equal(t, 2, itemCount)
}
