package main

import (
	"bookings/internal/config"
	"bookings/internal/logger"
	"bookings/internal/storage"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	logger.SetupLogger()

	slog.Info("application started")

	storage.NewPostgresDb(cfg)
	slog.Info("DB connected!")

}
