package order

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	sqlc "github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) Create(ctx context.Context, newOrder ordermodels.NewOrder) (orderID int64, err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "Create")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", newOrder.User),
		attribute.String("status", string(newOrder.Status)),
	)

	queries := sqlc.New(r.db.MasterDB())

	params := r.prepareCreateOrderParams(newOrder)

	tx, commit, rollback, err := transaction.CreateTx(ctx, r.db.MasterDB(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer rollback(ctx)

	queries = queries.WithTx(tx)

	start := time.Now()
	id, err := queries.CreateOrder(ctx, params)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("insert", err, duration)

	if err != nil {
		return 0, err
	}

	start = time.Now()
	err = queries.InsertOrderItems(ctx, r.convertToInsertOrderItemsParams(id, newOrder))
	duration = time.Since(start).Seconds()
	metric.RecordDBMetric("insert", err, duration)

	if err != nil {
		return 0, err
	}

	err = commit(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) prepareCreateOrderParams(newOrder ordermodels.NewOrder) sqlc.CreateOrderParams {
	return sqlc.CreateOrderParams{
		UserID: int32(newOrder.User),
		Status: sqlc.OrderStatus(newOrder.Status),
	}
}

func (r *Repository) convertToInsertOrderItemsParams(
	orderID int64,
	newOrder ordermodels.NewOrder,
) sqlc.InsertOrderItemsParams {

	skus := make([]int32, len(newOrder.Items))
	counts := make([]int32, len(newOrder.Items))
	for i, item := range newOrder.Items {
		skus[i] = int32(item.SKU)
		counts[i] = int32(item.Count)
	}

	return sqlc.InsertOrderItemsParams{
		OrderID: orderID,
		ItemSku: skus,
		Count:   counts,
	}
}
