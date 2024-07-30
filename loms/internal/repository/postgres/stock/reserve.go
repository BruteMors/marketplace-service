package stock

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	stockmodels "github.com/BruteMors/marketplace-service/loms/internal/models/stock"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/stock/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (r *Repository) Reserve(ctx context.Context, items []stockmodels.ReserveItem) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "Reserve")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	queries := sqlc.New(r.db.MasterDB())

	tx, commit, rollback, err := transaction.CreateTx(ctx, r.db.MasterDB(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer rollback(ctx)

	queries = queries.WithTx(tx)

	skus := make([]int32, 0, len(items))
	for _, item := range items {
		skus = append(skus, int32(item.SKU))
	}

	skuToCount := make(map[int32]int32, len(items))
	for _, item := range items {
		skuToCount[int32(item.SKU)] = int32(item.Count)
	}

	start := time.Now()
	itemsAvailable, err := queries.GetItemsAvailability(ctx, skus)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("select", err, duration)

	if err != nil {
		return err
	}

	for _, item := range itemsAvailable {
		if item.Available < skuToCount[item.Sku] {
			return repository.ErrInsufficientStock
		}
	}

	start = time.Now()
	err = queries.UpdateReservedItems(ctx, r.convertToReservedItemsParams(items))
	duration = time.Since(start).Seconds()
	metric.RecordDBMetric("update", err, duration)

	if err != nil {
		return err
	}

	start = time.Now()
	err = commit(ctx)
	duration = time.Since(start).Seconds()
	metric.RecordDBMetric("commit", err, duration)

	return err
}

func (r *Repository) convertToReservedItemsParams(items []stockmodels.ReserveItem) sqlc.UpdateReservedItemsParams {
	skus := make([]int32, len(items))
	counts := make([]int32, len(items))

	for i, item := range items {
		skus[i] = int32(item.SKU)
		counts[i] = int32(item.Count)
	}

	return sqlc.UpdateReservedItemsParams{
		Sku:   skus,
		Count: counts,
	}
}
