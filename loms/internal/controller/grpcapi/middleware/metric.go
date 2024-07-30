package middleware

import (
	"context"
	"regexp"
	"time"

	"github.com/BruteMors/marketplace-service/loms/internal/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	methodRegex = regexp.MustCompile(`/\w+/\d+`)
)

func RequestMetric(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	cleanedMethod := cleanURL(info.FullMethod)

	resp, err = handler(ctx, req)

	duration := time.Since(start).Seconds()
	metric.IncRequestCounter()
	statusCode := status.Code(err).String()
	metric.IncResponseCounter(statusCode, info.FullMethod, cleanedMethod)
	metric.ObserveResponseTime(statusCode, cleanedMethod, duration)

	return resp, err
}

func cleanURL(method string) string {
	return methodRegex.ReplaceAllString(method, "/{method}/{id}")
}
