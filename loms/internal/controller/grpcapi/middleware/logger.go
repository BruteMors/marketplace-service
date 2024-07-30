package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func RequestLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"received request",
		slog.String("method", info.FullMethod),
		slog.String("req", string(raw)),
	)

	return handler(ctx, req)
}
