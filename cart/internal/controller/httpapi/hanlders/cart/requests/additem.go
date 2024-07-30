package requests

type AddItem struct {
	UserID int64  `json:"-" validate:"required,gt=0"`
	SkuID  int64  `json:"-" validate:"required,gt=0"`
	Count  uint16 `json:"count" validate:"required,gt=0"`
}
