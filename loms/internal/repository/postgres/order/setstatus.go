package order

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	sqlc "github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/order/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status ordermodels.Status) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "SetStatus")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("orderID", orderID),
		attribute.String("status", string(status)),
	)

	queries := sqlc.New(r.db.MasterDB())

	params := r.prepareSetOrderStatusParams(orderID, status)

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	err = queries.SetOrderStatus(ctx, params)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("update", err, duration)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) prepareSetOrderStatusParams(orderID int64, status ordermodels.Status) sqlc.SetOrderStatusParams {
	return sqlc.SetOrderStatusParams{
		OrderID: orderID,
		Status:  sqlc.OrderStatus(status),
	}
}
