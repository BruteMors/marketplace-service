package order

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	grpcmodels "github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GRPCApi) OrderPay(
	ctx context.Context,
	in *grpcmodels.OrderPayRequest,
) (resp *emptypb.Empty, err error) {
	tracer := otel.Tracer("GRPCApi")
	propagator := otel.GetTextMapPropagator()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))

	var span trace.Span
	ctx, span = tracer.Start(ctx, "OrderPay")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("orderID", in.OrderId))

	err = g.orderService.OrderPay(ctx, in.OrderId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
