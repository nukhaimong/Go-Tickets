package bookings

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domain/event"
	"gotickets/internal/middlewares"
	"gotickets/internal/payment"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	bookingRepo := NewRepostory(db)
	eventRepo := event.NewRepository(db)
	stripeService := payment.NewStripeService(
		cfg.StripeSuccessURL,
		cfg.StripeCancelURL,
		cfg.StripeSecretKey,
	)
	bookingService := NewService(bookingRepo, eventRepo, stripeService)
	bookingHandler := NewHandler(bookingService)
	webhookHandler := payment.NewWebhookHandler(bookingService, cfg.StripeWebhookSecret)
	e.POST("/webhook/stripe", webhookHandler.HandleWebhook)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("api/v1/bookings", middlewares.AuthMiddleware(jwtService))

	api.POST("", bookingHandler.CreateBooking)
	api.GET("", bookingHandler.GetMyBookings)
	api.GET("/:id", bookingHandler.GetByID)
}
