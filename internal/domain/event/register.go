package event

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// dependency injection
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	api := e.Group("/api/v1/events")

	api.GET("", handler.GetEvents)
	api.GET("/:id", handler.GetEventById)
	api.POST("/create", handler.CreateEvent)
	api.PATCH("/:id", handler.UpdateEvent)
}
