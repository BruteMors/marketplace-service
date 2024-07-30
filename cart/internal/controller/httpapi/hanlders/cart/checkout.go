package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/requests"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/responses"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *HttpApi) Checkout(writer http.ResponseWriter, request *http.Request) (err error) {
	tr := otel.Tracer("httpApi")
	ctx, span := tr.Start(request.Context(), "Checkout")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	writer.Header().Set("Content-Type", "application/json")

	var req requests.Checkout

	req.UserID, err = strconv.ParseInt(request.PathValue("user_id"), 10, 64)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int64("userID", req.UserID))

	err = h.validator.Struct(req)
	if err != nil {
		return httpapi.ErrValidation
	}

	orderID, err := h.cartService.Checkout(ctx, req.UserID)
	if err != nil {
		return err
	}

	response := responses.Checkout{OrderID: orderID}

	responseBuf, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	_, err = writer.Write(responseBuf)
	if err != nil {
		return err
	}

	return nil
}
