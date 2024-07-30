package cart

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/metric"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) DeleteItemsByUserID(_ context.Context, userID int64) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	count := len(r.store[userID])
	delete(r.store, userID)

	return count, nil
}

func (p *RepositoryWithMetrics) DeleteItemsByUserID(ctx context.Context, userID int64) (count int, err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "DeleteItemsByUserID")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("userID", userID))

	start := time.Now()
	count, err = p.repo.DeleteItemsByUserID(ctx, userID)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	metric.IncDBRequestCounter("delete_items_by_user_id", status)
	metric.ObserveDBResponseTime("delete_items_by_user_id", status, duration)
	metric.SubInMemoryObjectCount(float64(count))

	return count, err
}
