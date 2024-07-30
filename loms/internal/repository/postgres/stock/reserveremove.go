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

func (r *Repository) ReserveRemove(ctx context.Context, item []stockmodels.ReserveItem) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "ReserveRemove")
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
	err = queries.ReserveRemove(ctx, r.convertToReserveRemoveItemsParams(item))
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("update", err, duration)

	return err
}

func (r *Repository) convertToReserveRemoveItemsParams(items []stockmodels.ReserveItem) sqlc.ReserveRemoveParams {
	skus := make([]int32, len(items))
	counts := make([]int32, len(items))

	for i, item := range items {
		skus[i] = int32(item.SKU)
		counts[i] = int32(item.Count)
	}

	return sqlc.ReserveRemoveParams{
		Sku:   skus,
		Count: counts,
	}
}
