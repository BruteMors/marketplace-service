package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GRPCServerConfig struct {
	host string
	port string
}

func NewGRPCServerConfig() (*GRPCServerConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if host == "" {
		return nil, errors.New("gRPC host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, errors.New("gRPC port not found")
	}

	return &GRPCServerConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *GRPCServerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
