package requests

type OrderCreate struct {
	User  int64
	Items []Item
}

type Item struct {
	SKU   uint32
	Count uint16
}
