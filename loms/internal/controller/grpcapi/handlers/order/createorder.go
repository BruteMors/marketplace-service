package order

import (
	"context"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/requests"
	grpcmodels "github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (g *GRPCApi) OrderCreate(
	ctx context.Context,
	in *grpcmodels.OrderCreateRequest,
) (resp *grpcmodels.OrderCreateResponse, err error) {
	tracer := otel.Tracer("GRPCApi")
	propagator := otel.GetTextMapPropagator()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))

	var span trace.Span
	ctx, span = tracer.Start(ctx, "OrderCreate")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	orderID, err := g.orderService.OrderCreate(ctx, repackOrderCreateRequest(in))
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	span.SetAttributes(
		attribute.Int64("orderID", orderID),
		attribute.Int64("userID", in.User),
	)
	resp = &grpcmodels.OrderCreateResponse{OrderId: orderID}
	return resp, nil
}

func repackOrderCreateRequest(in *grpcmodels.OrderCreateRequest) *requests.OrderCreate {
	items := make([]requests.Item, 0, len(in.Items))

	for _, i := range in.Items {
		items = append(items, requests.Item{
			SKU:   i.Sku,
			Count: uint16(i.Count),
		})
	}

	return &requests.OrderCreate{
		User:  in.User,
		Items: items,
	}
}
