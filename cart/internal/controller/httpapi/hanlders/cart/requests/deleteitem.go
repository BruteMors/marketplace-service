package requests

type DeleteItem struct {
	UserID int64 `json:"-" validate:"required,gt=0"`
	SkuID  int64 `json:"-" validate:"required,gt=0"`
}
