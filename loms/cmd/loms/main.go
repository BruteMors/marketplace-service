package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/BruteMors/marketplace-service/libs/logger"
	"github.com/BruteMors/marketplace-service/loms/internal/app"
)

func main() {
	handler := logger.NewCustomTextHandler(os.Stdout, "loms", nil)

	slog.SetDefault(slog.New(handler))

	ctx := context.Background()

	lomsApp, err := app.NewLoms(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = lomsApp.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
