package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func (s *ProductService) GetListSkus(ctx context.Context, startAfterSku, count int64) (response *models.ListSkusResponse, err error) {
	tr := otel.Tracer("productService")
	ctx, span := tr.Start(ctx, "GetListSkus")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(
		attribute.Int64("startAfterSku", startAfterSku),
		attribute.Int64("count", count),
	)

	url := fmt.Sprintf("%s/list_skus", s.address)

	reqBody := models.ListSkusRequest{
		Token:         s.accessToken,
		StartAfterSku: startAfterSku,
		Count:         count,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

	httpReq = httpReq.WithContext(ctx)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		errBodyClose := Body.Close()
		if errBodyClose != nil {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"error with Body.Close()",
				slog.String("method", "ProductService.GetListSkus"),
				slog.String("error", errBodyClose.Error()),
			)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var listSkusResponse models.ListSkusResponse
	if err := json.NewDecoder(resp.Body).Decode(&listSkusResponse); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &listSkusResponse, nil
}
