package cart

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/requests"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/responses"
	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *HttpApi) GetCart(writer http.ResponseWriter, request *http.Request) (err error) {
	tr := otel.Tracer("httpApi")
	ctx, span := tr.Start(request.Context(), "GetCart")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	writer.Header().Set("Content-Type", "application/json")

	var req requests.GetCart

	req.UserID, err = strconv.ParseInt(request.PathValue("user_id"), 10, 64)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int64("userID", req.UserID))

	err = h.validator.Struct(req)
	if err != nil {
		return httpapi.ErrValidation
	}

	cart, err := h.cartService.GetCart(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			writer.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}

	items := make([]responses.Item, 0, len(cart.Items))

	for _, item := range cart.Items {
		items = append(items, responses.Item{
			SkuID: item.SkuID,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
	}

	response := responses.GetCart{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}

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
