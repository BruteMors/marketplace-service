package outbox

import (
	"context"
	"time"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	"github.com/BruteMors/marketplace-service/loms/internal/repository/postgres/outbox/sqlc"
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/transaction"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *Repository) MarkOrderStatusChangedEventAsSend(ctx context.Context, eventID int64) (err error) {
	tr := otel.Tracer("repository")
	ctx, span := tr.Start(ctx, "MarkOrderStatusChangedEventAsSend")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("eventID", eventID),
	)

	queries := sqlc.New(r.db.MasterDB())

	tx, found := transaction.CheckTx(ctx)
	if found {
		queries = queries.WithTx(tx)
	}

	start := time.Now()
	err = queries.MarkOrderStatusChangedEventAsSend(ctx, eventID)
	duration := time.Since(start).Seconds()
	metric.RecordDBMetric("update", err, duration)

	return err
}
