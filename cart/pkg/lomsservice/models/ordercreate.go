package models

type OrderCreate struct {
	User  int64
	Items []OrderItem
}

type OrderItem struct {
	SkuID uint32
	Count uint32
}
