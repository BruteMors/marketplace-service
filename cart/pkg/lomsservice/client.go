package lomsservice

import (
	"errors"
	"net"
	"os"

	"github.com/BruteMors/marketplace-service/cart/pkg/api/grpc/loms/v1"
	"github.com/BruteMors/marketplace-service/cart/pkg/lomsservice/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcLomsServiceHostEnvName = "GRPC_LOMS_SERVICE_HOST"
	grpcLomsServicePortEnvName = "GRPC_LOMS_SERVICE_PORT"
)

type Client struct {
	orderClient loms.OrdersClient
	stockClient loms.StockClient
}

func NewClient() (*Client, error) {
	host := os.Getenv(grpcLomsServiceHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc loms host not found")
	}

	port := os.Getenv(grpcLomsServicePortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc loms port not found")
	}

	address := net.JoinHostPort(host, port)

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.SetTraceID),
	)
	if err != nil {
		return nil, err
	}

	orderClient := loms.NewOrdersClient(conn)
	stockClient := loms.NewStockClient(conn)

	return &Client{orderClient: orderClient, stockClient: stockClient}, nil
}
