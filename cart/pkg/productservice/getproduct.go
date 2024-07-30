package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/BruteMors/marketplace-service/cart/pkg/errorgroup"
	"github.com/BruteMors/marketplace-service/cart/pkg/productservice/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/time/rate"
)

func (s *ProductService) GetProduct(ctx context.Context, sku int64) (response *models.GetProductResponse, err error) {
	tr := otel.Tracer("productService")
	ctx, span := tr.Start(ctx, "GetProduct")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("sku", sku))

	url := fmt.Sprintf("%s/get_product", s.address)

	reqBody := models.GetProductRequest{
		Token: s.accessToken,
		Sku:   sku,
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
				slog.String("method", "ProductService.GetProduct"),
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

	var productResponse models.GetProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResponse); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	productResponse.Sku = sku
	return &productResponse, nil
}

func (s *ProductService) GetProducts(ctx context.Context, skus []int64) (responses []models.GetProductResponse, err error) {
	tr := otel.Tracer("productService")
	ctx, span := tr.Start(ctx, "GetProducts")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64Slice("skus", skus))

	responses = make([]models.GetProductResponse, len(skus))

	limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(s.getProductRPSLimit)), 1)

	eg := errorgroup.NewErrGroup(ctx)

	for i, sku := range skus {
		i, sku := i, sku
		eg.Go(func() error {
			if err := limiter.Wait(eg.Context()); err != nil {
				return err
			}

			product, err := s.GetProduct(eg.Context(), sku)
			if err != nil {
				return err
			}

			responses[i].Name = product.Name
			responses[i].Price = product.Price
			responses[i].Sku = product.Sku

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return nil, err
	}

	return responses, nil
}
