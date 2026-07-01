package bookings

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/event"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	bookingRepo := NewRepostory(db)
	eventRepo := event.NewRepository(db)
	bookingService := NewService(bookingRepo, eventRepo)
	bookingHandler := NewHandler(bookingService)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("api/v1/bookings", middlewares.AuthMiddleware(jwtService))

	api.POST("", bookingHandler.CreateBooking)
	api.GET("", bookingHandler.GetMyBookings)
}
