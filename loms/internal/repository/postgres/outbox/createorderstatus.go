package outbox

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/outbox/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) CreateOrderStatusChangedEvent(ctx context.Context, orderID int64, status ordermodels.Status) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "CreateOrderStatusChangedEvent")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("orderID", orderID),
		attribute.String("status", string(status)),
	)

	queries := sqlc.New(r.db.MasterDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	err = queries.CreateOrderStatusChangedEvent(ctx, r.convertToCreateOrderStatusChangedEventParams(orderID, status))
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("insert", err, duration)

	return err
}

func (r *Repository) convertToCreateOrderStatusChangedEventParams(
	orderID int64,
	status ordermodels.Status,
) sqlc.CreateOrderStatusChangedEventParams {
	return sqlc.CreateOrderStatusChangedEventParams{
		OrderID: orderID,
		Status:  sqlc.OrderStatus(status),
	}
}
