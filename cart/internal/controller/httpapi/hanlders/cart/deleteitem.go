package cart

import (
	"net/http"
	"strconv"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/requests"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *HttpApi) DeleteItem(writer http.ResponseWriter, request *http.Request) (err error) {
	tr := otel.Tracer("httpApi")
	ctx, span := tr.Start(request.Context(), "DeleteItem")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	writer.Header().Set("Content-Type", "application/json")

	var req requests.DeleteItem

	req.UserID, err = strconv.ParseInt(request.PathValue("user_id"), 10, 64)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int64("userID", req.UserID))

	req.SkuID, err = strconv.ParseInt(request.PathValue("sku_id"), 10, 64)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int64("skuID", req.SkuID))

	err = h.validator.Struct(req)
	if err != nil {
		return httpapi.ErrValidation
	}

	err = h.cartService.DeleteItem(ctx, req.UserID, req.SkuID)
	if err != nil {
		return err
	}

	writer.WriteHeader(http.StatusNoContent)

	return nil
}
