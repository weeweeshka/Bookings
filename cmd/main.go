package main

import (
	"bookings/internal/config"
	"bookings/internal/logger"
	"bookings/internal/router"
	"bookings/internal/storage"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	logger.SetupLogger()

	slog.Info("application started")

	postgres, err := storage.NewPostgresDb(cfg)
	if err != nil {
		slog.Error("failed to init storage", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	slog.Info("DB connected!")

	_ = postgres

	router := router.SetupRouter()
	router.Run(":8080")

}
