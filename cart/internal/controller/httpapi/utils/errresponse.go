package utils

import (
	"encoding/json"
	"net/http"

	"github.com/BruteMors/marketplace-service/cart/internal/controller/httpapi"
)

func WriteErrResponse(writer http.ResponseWriter, statusCode int, err error) error {
	writer.WriteHeader(statusCode)
	var errorHandler httpapi.Error
	errorHandler.Message = httpapi.ErrInvalidSku.Error()

	buf, err := json.Marshal(errorHandler)
	if err != nil {
		return err
	}

	_, err = writer.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
