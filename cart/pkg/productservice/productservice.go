package productservice

import (
	"github.com/BruteMors/marketplace-service/cart/pkg/httpclient"
)

type ProductService struct {
	client             *httpclient.HttpClient
	address            string
	accessToken        string
	getProductRPSLimit int
}

func New(client *httpclient.HttpClient, address string, accessToken string, getProductRPSLimit int) *ProductService {
	return &ProductService{
		client:             client,
		address:            address,
		accessToken:        accessToken,
		getProductRPSLimit: getProductRPSLimit,
	}
}
