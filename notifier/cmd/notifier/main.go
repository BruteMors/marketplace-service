package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/BruteMors/marketplace-service/libs/logger"
	"github.com/BruteMors/marketplace-service/notifier/internal/app"
)

func main() {
	handler := logger.NewCustomTextHandler(os.Stdout, "notifier", nil)

	slog.SetDefault(slog.New(handler))

	ctx := context.Background()

	notifierApp, err := app.NewNotifier(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = notifierApp.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
