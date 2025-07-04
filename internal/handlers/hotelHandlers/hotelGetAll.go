package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetHotels interface {
	GetAllHotels() (string, error)
}

func GetAllHotelHandler(log *slog.Logger, getHotels GetHotels) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.hotelHandlers.GetAllHotelHandler"

		slog.With(slog.String("op", op))

		data, err := getHotels.GetAllHotels()
		if err != nil {
			slog.Info("failed to get hotels")

			c.JSON(http.StatusBadRequest, gin.H{"err": err})

			return
		}

		c.Data(200, "application/json", []byte(data))
		slog.Info("succssesful")
	}
}
