package models

type Error struct {
	Message string
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
	ErrProductNotFound = NewError("product not found")
	ErrCartNotFound    = NewError("cart not found")
	ErrStocksNotEnough = NewError("stocks not enough")
)
