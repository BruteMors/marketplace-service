package httpapi

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(message string) Error {
	return Error{
		Message: message,
	}
}

var (
	ErrValidation = NewError("validation error")
	ErrInvalidSku = NewError("invalid sku")
)
