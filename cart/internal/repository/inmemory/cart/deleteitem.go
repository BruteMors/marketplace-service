package cart

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/metric"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) DeleteItem(_ context.Context, userID int64, skuID int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.store[userID]; !ok {
		return nil
	}

	delete(r.store[userID], skuID)

	return nil
}

func (p *RepositoryWithMetrics) DeleteItem(ctx context.Context, userID int64, skuID int64) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "DeleteItem")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", userID),
		attribute.Int64("skuID", skuID),
	)

	start := time.Now()
	err = p.repo.DeleteItem(ctx, userID, skuID)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	metric.IncDBRequestCounter("delete", status)
	metric.ObserveDBResponseTime("delete", status, duration)
	metric.SubInMemoryObjectCount(1)

	return err
}
