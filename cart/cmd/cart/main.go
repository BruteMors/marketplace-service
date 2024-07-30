package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/BruteMors/marketplace-service/cart/internal/app"
	"github.com/BruteMors/marketplace-service/libs/logger"
)

func main() {
	handler := logger.NewCustomTextHandler(os.Stdout, "cart", nil)

	slog.SetDefault(slog.New(handler))

	ctx := context.Background()

	cartApp, err := app.NewCart(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = cartApp.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
