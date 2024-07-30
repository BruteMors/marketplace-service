package cart

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/metric"
	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/cart/internal/repository"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) GetCart(_ context.Context, userID int64) ([]models.ItemCount, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.store[userID]; !ok {
		return nil, repository.ErrCartNotFound
	}

	if len(r.store[userID]) == 0 {
		return nil, repository.ErrCartEmpty
	}

	items := make([]models.ItemCount, 0, len(r.store[userID]))

	for skuID, count := range r.store[userID] {
		items = append(items, models.ItemCount{
			SkuID: skuID,
			Count: count,
		})
	}

	return items, nil
}

func (p *RepositoryWithMetrics) GetCart(ctx context.Context, userID int64) (items []models.ItemCount, err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "GetCart")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("userID", userID))

	start := time.Now()
	items, err = p.repo.GetCart(ctx, userID)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	metric.IncDBRequestCounter("get_cart", status)
	metric.ObserveDBResponseTime("get_cart", status, duration)

	return items, err
}
