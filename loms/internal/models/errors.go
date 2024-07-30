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
	ErrSKUNotFound   = NewError("sku not found")
	ErrOrderNotFound = NewError("order not found")
)
