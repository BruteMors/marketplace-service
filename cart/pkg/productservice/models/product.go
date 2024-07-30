package models

type GetProductRequest struct {
	Token string `json:"token"`
	Sku   int64  `json:"sku"`
}

type GetProductResponse struct {
	Sku   int64  `json:"-"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type ListSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku int64  `json:"startAfterSku"`
	Count         int64  `json:"count"`
}

type ListSkusResponse struct {
	Skus []int64 `json:"skus"`
}
