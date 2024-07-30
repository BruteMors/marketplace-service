package requests

type GetCart struct {
	UserID int64 `json:"-" validate:"required,gt=0"`
}
