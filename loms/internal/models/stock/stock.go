package stock

type Item struct {
	SKU        uint32
	TotalCount uint64
	Reserved   uint64
}

type ReserveItem struct {
	SKU   uint32
	Count uint16
}
