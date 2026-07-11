package bookings

import (
	"errors"
	"gotickets/internal/domain/bookings/dto"
	"gotickets/internal/domain/event"
	httpresponse "gotickets/internal/httpResponse"
	"gotickets/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{service: s}
}

func bookingErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrBookingNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Booking not found",
		})
	}

	if errors.Is(err, event.ErrEventNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Event not found",
		})
	}

	if errors.Is(err, ErrNotEnoughTickets) {
		return c.JSON(http.StatusConflict, httpresponse.Error{
			Code:    http.StatusConflict,
			Message: "Not enough tickets available",
		})
	}

	if errors.Is(err, ErrBookingAlreadyCancelled) {
		return c.JSON(http.StatusConflict, httpresponse.Error{
			Code:    http.StatusConflict,
			Message: "Booking is already cancelled",
		})
	}

	if errors.Is(err, ErrForbiddenBookingAccess) {
		return c.JSON(http.StatusForbidden, httpresponse.Error{
			Code:    http.StatusForbidden,
			Message: "You do not own this booking",
		})
	}

	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateBooking(c *echo.Context) error {
	userId, ok := utils.GetCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorize",
		})
	}

	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateBooking(userId, req)
	if err != nil {
		return bookingErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetMyBookings(c *echo.Context) error {
	userId, ok := utils.GetCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorize",
		})
	}
	bookings, err := h.service.GetMyBookings(userId)
	if err != nil {
		return bookingErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, bookings)
}

func(h *handler) GetByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Ivalid event id",
			Details: err.Error(),
		})
	}
	event, err := h.service.GetByID(uint(id))
	if err != nil {
		return bookingErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, event)
}
