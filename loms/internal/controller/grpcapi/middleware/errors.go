package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandler(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		if _, ok := status.FromError(err); !ok {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"request failed with unexpected error",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
			)
			return nil, status.Error(codes.Internal, "internal server error")
		}

		return nil, err
	}

	return resp, nil
}
