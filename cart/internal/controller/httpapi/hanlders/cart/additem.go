package cart

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/requests"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/hanlders/cart/responses"
	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi/utils"
	"github.com/BruteMors/marketplace-service/cart/internal/models"
	"github.com/BruteMors/marketplace-service/libs/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *HttpApi) AddItem(writer http.ResponseWriter, request *http.Request) (err error) {
	tr := otel.Tracer("httpApi")
	ctx, span := tr.Start(request.Context(), "AddItem")
	defer func() {
		tracing.RecordSpanError(span, err)
		span.End()
	}()

	writer.Header().Set("Content-Type", "application/json")

	var req requests.AddItem

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

	buf, err := io.ReadAll(request.Body)
	defer func(Body io.ReadCloser) {
		errBodyClose := Body.Close()
		if errBodyClose != nil {
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"error with Body.Close()",
				slog.String("method", "HttpApi.AddItem"),
				slog.String("error", errBodyClose.Error()),
			)
		}
	}(request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &req)
	if err != nil {
		return err
	}

	err = h.validator.Struct(req)
	if err != nil {
		return httpapi.ErrValidation
	}

	span.SetAttributes(attribute.Int("count", int(req.Count)))

	err = h.cartService.AddItem(ctx, req.UserID, req.SkuID, req.Count)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			err = utils.WriteErrResponse(writer, http.StatusPreconditionFailed, httpapi.ErrInvalidSku)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	response := responses.AddItem{}

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
