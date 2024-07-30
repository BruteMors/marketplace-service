package stock

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	sqlc "github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) GetBySKU(ctx context.Context, skuID uint32) (item stockmodels.Item, err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "GetBySKU")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("skuID", int64(skuID)))

	queries := sqlc.New(r.db.ReplicaDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	dbOrder, err := queries.GetBySKU(ctx, int32(skuID))
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("select", err, duration)

	if err != nil {
		return stockmodels.Item{}, err
	}

	item = r.convertGetBySKURowToItem(dbOrder)
	return item, nil
}

func (r *Repository) convertGetBySKURowToItem(row sqlc.GetBySKURow) stockmodels.Item {
	return stockmodels.Item{
		SKU:        uint32(row.Sku),
		TotalCount: uint64(row.TotalCount),
		Reserved:   uint64(row.Reserved),
	}
}
