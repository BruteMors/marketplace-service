package httpclient

import (
	"net/http"
	"time"

	"github.com/BruteMors/marketplace-service/cart/internal/metric"
)

type HttpClient struct {
	http.Client
}

type retryTransport struct {
	transport  http.RoundTripper
	retries    int
	retryDelay time.Duration
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= rt.retries; i++ {
		start := time.Now()
		resp, err = rt.transport.RoundTrip(req)
		duration := time.Since(start).Seconds()

		status := "error"
		if resp != nil {
			status = http.StatusText(resp.StatusCode)
		}

		metric.IncExternalRequestCounter(status, req.URL.Path)
		metric.ObserveExternalResponseTime(status, req.URL.Path, duration)

		if err == nil && resp != nil && resp.StatusCode != 420 && resp.StatusCode != 429 {
			return resp, nil
		}

		time.Sleep(rt.retryDelay)
	}

	return resp, err
}

func New(timeout time.Duration, retries int, retryDelay time.Duration) *HttpClient {
	httpClient := http.Client{
		Timeout: timeout,
		Transport: &retryTransport{
			transport:  http.DefaultTransport,
			retries:    retries,
			retryDelay: retryDelay,
		},
	}

	return &HttpClient{
		httpClient,
	}
}
