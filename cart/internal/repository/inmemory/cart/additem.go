package cart

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/metric"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) Add(_ context.Context, userID int64, skuID int64, count uint16) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.store[userID]; !ok {
		r.store[userID] = make(map[int64]uint16)
	}

	r.store[userID][skuID] += count

	return nil
}

func (p *RepositoryWithMetrics) Add(ctx context.Context, userID int64, skuID int64, count uint16) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "Add")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("userID", userID),
		attribute.Int64("skuID", skuID),
		attribute.Int64("count", int64(count)),
	)

	start := time.Now()
	err = p.repo.Add(ctx, userID, skuID, count)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	metric.IncDBRequestCounter("add", status)
	metric.ObserveDBResponseTime("add", status, duration)
	metric.AddInMemoryObjectCount(1)

	return err
}
