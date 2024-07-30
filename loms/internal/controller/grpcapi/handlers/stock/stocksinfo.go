package stock

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	grpcmodels "github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (g *GRPCApi) StocksInfo(
	ctx context.Context,
	in *grpcmodels.StocksInfoRequest,
) (resp *grpcmodels.StocksInfoResponse, err error) {
	tracer := otel.Tracer("GRPCApi")
	var span trace.Span
	ctx, span = tracer.Start(ctx, "StocksInfo")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("sku", int64(in.Sku)))

	stockCount, err := g.stockService.StocksInfo(ctx, in.Sku)
	if err != nil {
		return nil, err
	}

	return &grpcmodels.StocksInfoResponse{Count: stockCount}, nil
}
