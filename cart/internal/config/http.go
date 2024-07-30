package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName             = "HTTP_HOST"
	httpPortEnvName             = "HTTP_PORT"
	httpClintTimeoutEnvName     = "HTTP_CLIENT_TIMEOUT"
	httpClientRetriesEnvName    = "HTTP_CLIENT_RETRIES"
	httpClientRetryDelayEnvName = "HTTP_CLIENT_RETRY_DELAY"
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

type HTTPClientConfig struct {
	Timeout    string
	Retries    string
	RetryDelay string
}

func NewHTTPClientConfig() (*HTTPClientConfig, error) {
	timeout := os.Getenv(httpClintTimeoutEnvName)
	if len(timeout) == 0 {
		return nil, errors.New("http client timeout not found")
	}

	retries := os.Getenv(httpClientRetriesEnvName)
	if len(retries) == 0 {
		return nil, errors.New("http client retries not found")
	}

	retryDelay := os.Getenv(httpClientRetryDelayEnvName)
	if len(retryDelay) == 0 {
		return nil, errors.New("http client retry delay not found")
	}

	return &HTTPClientConfig{
		Timeout:    timeout,
		RetryDelay: retryDelay,
		Retries:    retries,
	}, nil
}
