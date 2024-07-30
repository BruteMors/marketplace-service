package requests

type Checkout struct {
	UserID int64 `json:"-" validate:"required,gt=0"`
}
