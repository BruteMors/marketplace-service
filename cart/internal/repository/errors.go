package repository

import "errors"

var (
	ErrCartNotFound = errors.New("cart not found")
	ErrCartEmpty    = errors.New("cart is empty")
)
