package outbox

import (
	"context"
	"errors"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	ordermodels "github.com/BruteMors/marketplace-service/loms/internal/models/order"
	"github.com/BruteMors/marketplace-service/loms/internal/repository"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/outbox/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (r *Repository) FetchNextOrderStatusChangedEvent(ctx context.Context) (
	event ordermodels.StatusChangedEvent,
	err error,
) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "FetchNextOrderStatusChangedEvent")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	queries := sqlc.New(r.db.ReplicaDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	dbEvent, err := queries.FetchNextOrderStatusChangedEvent(ctx)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("select", err, duration)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ordermodels.StatusChangedEvent{}, repository.ErrNoElements
		}
		return ordermodels.StatusChangedEvent{}, err
	}

	event = ordermodels.StatusChangedEvent{
		ID:      dbEvent.ID,
		OrderID: dbEvent.OrderID,
		Status:  ordermodels.Status(dbEvent.Status),
		At:      dbEvent.At.Time,
	}

	return event, nil
}
