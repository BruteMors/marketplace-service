package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	productServiceAddressEnvName            = "PRODUCT_SERVICE_ADDRESS"
	productServiceTokenEnvName              = "PRODUCT_SERVICE_TOKEN"
	productServiceGetProductRPSLimitEnvName = "PRODUCT_SERVICE_GET_PRODUCT_RPS_LIMIT"
)

type ProductServiceConfig struct {
	Address            string
	Token              string
	GetProductRPSLimit int
}

func NewProductServiceConfig() (*ProductServiceConfig, error) {
	address := os.Getenv(productServiceAddressEnvName)
	if address == "" {
		return nil, errors.New("product service address is not set")
	}

	token := os.Getenv(productServiceTokenEnvName)
	if token == "" {
		return nil, errors.New("product service token is not set")
	}

	getProductRPSLimit := os.Getenv(productServiceGetProductRPSLimitEnvName)
	if getProductRPSLimit == "" {
		return nil, errors.New("product service get product rps limit is not set")
	}

	limit, err := strconv.Atoi(getProductRPSLimit)
	if err != nil {
		return nil, err
	}

	return &ProductServiceConfig{
		Address:            address,
		Token:              token,
		GetProductRPSLimit: limit,
	}, nil
}
