package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type HTTPServerConfig struct {
	host string
	port string
}

func NewHTTPServerConfig() (*HTTPServerConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &HTTPServerConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *HTTPServerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
