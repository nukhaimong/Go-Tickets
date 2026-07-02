package server

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/bookings"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally return the error to let each route control the status code.
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()
	db.AutoMigrate(&user.User{}, &event.Event{}, &bookings.Booking{})

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})
	e.GET("/jekono", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, "Hello From Jekono.com!")
	})

	// user routes
	user.RegisterRoutes(e, db, cfg)
	//event routes
	event.RegisterRoutes(e, db)
	// bookings routes
	bookings.RegisterRoutes(e, db, cfg)
	// port
	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
