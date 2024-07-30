package repository

import "errors"

var (
	ErrSKUNotFound       = errors.New("sku not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOrderNotFound     = errors.New("order not found")
	ErrNoElements        = errors.New("no elements")
)
