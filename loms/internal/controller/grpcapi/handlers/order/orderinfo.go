package order

import (
	"context"
	"errors"

	"github.com/BruteMors/marketplace-service/libs/tracing"
	"github.com/BruteMors/marketplace-service/loms/internal/controller/grpcapi/utils"
	"github.com/BruteMors/marketplace-service/loms/internal/models"
	"github.com/BruteMors/marketplace-service/loms/internal/models/order/responses"
	grpcmodels "github.com/BruteMors/marketplace-service/loms/pkg/api/grpc/loms/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (g *GRPCApi) OrderInfo(
	ctx context.Context,
	in *grpcmodels.OrderInfoRequest,
) (resp *grpcmodels.OrderInfoResponse, err error) {
	tracer := otel.Tracer("GRPCApi")
	propagator := otel.GetTextMapPropagator()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))

	var span trace.Span
	ctx, span = tracer.Start(ctx, "OrderInfo")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	span.SetAttributes(attribute.Int64("orderID", in.OrderId))

	orderInfo, err := g.orderService.OrderInfo(ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	response, err := repackOrderInfoResponse(&orderInfo)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func repackOrderInfoResponse(orderInfo *responses.OrderInfo) (*grpcmodels.OrderInfoResponse, error) {
	items := make([]*grpcmodels.OrderItem, 0, len(orderInfo.Items))

	for _, i := range orderInfo.Items {
		items = append(items, &grpcmodels.OrderItem{
			Sku:   i.SKU,
			Count: uint32(i.Count),
		})
	}

	status, err := utils.StringToGRPCOrderStatus(orderInfo.Status)
	if err != nil {
		return nil, err
	}

	return &grpcmodels.OrderInfoResponse{
		Status: status,
		User:   orderInfo.User,
		Items:  items,
	}, nil
}
