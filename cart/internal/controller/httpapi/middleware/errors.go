package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
	"github.com/BruteMors/marketplace-service/cart/internal/models"
)

type ErrorWrapper func(writer http.ResponseWriter, request *http.Request) error

func (s ErrorWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := s(writer, request); err != nil {
		ctx := request.Context()
		if !errors.As(err, &httpapi.Error{}) && !errors.As(err, &models.Error{}) {
			writer.WriteHeader(http.StatusInternalServerError)
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"request failed with unexpected error",
				slog.String("method", request.Method),
				slog.String("URL", request.URL.String()),
				slog.String("remote_addr", request.RemoteAddr),
				slog.String("error", err.Error()),
			)
			return
		}

		writer.WriteHeader(http.StatusBadRequest)

		var errorHandler httpapi.Error
		errorHandler.Message = err.Error()

		buf, err := json.Marshal(errorHandler)
		if err != nil {
			return
		}

		_, err = writer.Write(buf)
		if err != nil {
			return
		}
	}
}
