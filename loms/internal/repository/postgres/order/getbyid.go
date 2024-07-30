package order

import (
	"context"
	"fmt"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (order ordermodels.Order, err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "GetByID")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("orderID", orderID),
	)

	queries := sqlc.New(r.db.ReplicaDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	dbOrder, err := queries.GetByID(ctx, orderID)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("select", err, duration)

	if err != nil {
		return ordermodels.Order{}, err
	}

	order, err = r.convertGetByIDRowToOrder(dbOrder)
	if err != nil {
		return ordermodels.Order{}, err
	}

	return order, nil
}

func (r *Repository) convertGetByIDRowToOrder(dbOrder sqlc.GetByIDRow) (ordermodels.Order, error) {
	skus := dbOrder.Skus
	counts := dbOrder.Counts

	if len(skus) != len(counts) {
		return ordermodels.Order{}, fmt.Errorf("mismatched lengths of Skus and Counts")
	}

	items := make([]ordermodels.Item, len(skus))
	for i, sku := range skus {
		items[i] = ordermodels.Item{
			SKU:   uint32(sku),
			Count: uint16(counts[i]),
		}
	}

	createdAt, updatedAt := dbOrder.CreatedAt.Time, dbOrder.UpdatedAt.Time

	return ordermodels.Order{
		ID:        dbOrder.ID,
		Status:    ordermodels.Status(dbOrder.Status),
		UserID:    int64(dbOrder.UserID),
		Items:     items,
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}
