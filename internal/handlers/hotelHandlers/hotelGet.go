package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetHotel interface {
	GetHotel(id int) (string, error)
}

func GetHotelHandler(log slog.Logger, getHotel GetHotel) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.hotelHandlers.GetHotelHandler"

		slog := slog.With(slog.String("op", op))

		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		data, err := getHotel.GetHotel(id)
		if err != nil {
			slog.Info("failed to get hotel")

			c.JSON(http.StatusBadRequest, gin.H{"err": err})

			return
		}

		c.JSON(200, data)
	}
}
