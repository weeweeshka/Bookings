package main

import (
	"bookings/internal/config"
	"bookings/internal/logger"
)

func main() {
	config.MustLoad()
	logger.SetupLogger()

}
