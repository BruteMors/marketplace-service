package stock

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) ReserveCancel(ctx context.Context, item []stockmodels.ReserveItem) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "ReserveCancel")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	queries := sqlc.New(r.db.MasterDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	span.SetAttributes(
		attribute.Int("item_count", len(item)),
	)

	start := time.Now()
	err = queries.ReserveCancel(ctx, r.convertToReserveCancelItemsParams(item))
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("update", err, duration)

	return err
}

func (r *Repository) convertToReserveCancelItemsParams(items []stockmodels.ReserveItem) sqlc.ReserveCancelParams {
	skus := make([]int32, len(items))
	counts := make([]int32, len(items))

	for i, item := range items {
		skus[i] = int32(item.SKU)
		counts[i] = int32(item.Count)
	}

	return sqlc.ReserveCancelParams{
		Sku:   skus,
		Count: counts,
	}
}
