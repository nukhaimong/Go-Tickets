package event

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middlewares"
	"log"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// dependency injection
	// Initialize Cloudinary
	cloudinaryService, err := config.NewCloudinaryService()
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
	repo := NewRepository(db)
	service := NewService(repo, cloudinaryService)
	handler := NewHandler(service)
	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("/api/v1/events")

	api.GET("", handler.GetEvents)
	api.GET("/:id", handler.GetEventById)
	api.POST("/create", handler.CreateEvent, middlewares.AuthMiddleware(jwtService))
	api.PATCH("/:id", handler.UpdateEvent, middlewares.AuthMiddleware(jwtService))
}
