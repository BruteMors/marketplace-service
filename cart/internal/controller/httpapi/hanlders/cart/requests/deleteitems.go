package requests

type DeleteItems struct {
	UserID int64 `json:"-" validate:"required,gt=0"`
}
