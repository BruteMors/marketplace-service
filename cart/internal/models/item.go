package models

type Item struct {
	Name  string
	Price uint32
	ItemCount
}

type ItemCount struct {
	SkuID int64
	Count uint16
}
