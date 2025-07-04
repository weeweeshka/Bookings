package main

import (
	"bookings/internal/config"
	"bookings/internal/logger"
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

	postgres.CreateHotel("Россия", "Нальчик", "Веселый сыпыс", 5)
	postgres.CreateHotelRoom(1, 2, true, true, true, false)
	postgres.CreateVisitor(1, 1, "Сыпыс", "Сыпысвич", 56)

}
