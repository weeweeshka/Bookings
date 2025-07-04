package handlers

import (
	"bookings/internal/logger"
	"bookings/internal/models"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateHotel interface {
	CreateHotel(country string, city string, hotelName string, stars int) error
}

func PostHotelHandler(log *slog.Logger, createHotel CreateHotel) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.hotelHandlers.PostHotelHandler"
		var hotel models.Hotel

		slog := slog.With(
			slog.String("op", op),
		)

		err := c.ShouldBindJSON(&hotel)
		if errors.Is(err, io.EOF) {
			slog.Info("request body is empty")

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			slog.Info("failed to decode request body", logger.Err(err))

			return
		}

		slog.Info("request body decoded")

		err = createHotel.CreateHotel(hotel.Country, hotel.City, hotel.HotelName, hotel.Stars)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})

			slog.Info("failed to create hotel")

			return
		}

		c.JSON(200, hotel)
		slog.Info("HOTEL CREATED")
	}
}
