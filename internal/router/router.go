package router

import (
	"bookings/internal/config"
	handlers "bookings/internal/handlers/hotelHandlers"
	"bookings/internal/storage"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	cfg := config.MustLoad()
	postgres, _ := storage.NewPostgresDb(cfg)

	groupHotels := r.Group("/hotel")
	groupHotels.POST("/", handlers.PostHotelHandler(slog.Default(), postgres))
	groupHotels.GET("/", handlers.GetAllHotelHandler(slog.Default(), postgres))
	groupHotels.GET("/:id", handlers.GetHotelHandler(*slog.Default(), postgres))
	// groupHotels.DELETE("/:id", )
	// groupHotels.PUT("/:id", )
	return r
}
