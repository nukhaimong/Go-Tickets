package event

import (
	"errors"
	"gotickets/internal/domain/event/dto"
	httpresponse "gotickets/internal/httpResponse"
	"gotickets/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{
		service: s,
	}
}

func eventErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrEventNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Event not found",
		})
	}
	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateEvent(c *echo.Context) error {
	var req dto.CreateRequest

	// if err := c.Bind(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, httpresponse.Error{
	// 		Code:    http.StatusBadRequest,
	// 		Message: "Invalid request payload",
	// 		Details: err.Error(),
	// 	})
	// }

	// get user id from middleware
	userId, ok := utils.GetCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorize",
		})
	}
	//Bind form data
	req.Title = c.FormValue("title")
	req.Description = c.FormValue("description")
	req.Location = c.FormValue("location")
	startsAtStr := c.FormValue("starts_at")

	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to load timezone",
			Details: err.Error(),
		})
	}
	startsAt, err := time.ParseInLocation("2006-01-02 15:04:05", startsAtStr, loc)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid date format",
			Details: err.Error(),
		})
	}
	req.StartsAt = startsAt

	// Bind int and float
	totalTickets, _ := strconv.Atoi(c.FormValue("total_tickets"))
	req.TotalTickets = totalTickets

	priceStr := c.FormValue("price")

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid price value",
			Details: err.Error(),
		})

	}
	req.Price = price

	// bind the file
	fileHeader, err := c.FormFile("photo")
	if err == nil {
		req.Photo = fileHeader
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateEvent(&req, userId)
	if err != nil {
		return eventErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetEvents(c *echo.Context) error {
	events, err := h.service.GetEvents()
	if err != nil {
		return eventErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, events)
}

func (h *handler) GetEventById(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Ivalid event id",
			Details: err.Error(),
		})
	}
	event, err := h.service.GetEventById(uint(id))
	if err != nil {
		return eventErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, event)
}

func (h *handler) UpdateEvent(c *echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid event id",
			Details: err.Error(),
		})
	}
	userId, ok := utils.GetCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorize",
		})
	}

	var req dto.UpdateRequest

	// if err := c.Bind(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, httpresponse.Error{
	// 		Code:    http.StatusBadRequest,
	// 		Message: "Invalid request payload",
	// 		Details: err.Error(),
	// 	})
	// }

	title := c.FormValue("title")
	if title != "" {
		req.Title = &title
	}
	description := c.FormValue("description")
	if description != "" {
		req.Description = &description
	}
	location := c.FormValue("location")
	if location != "" {
		req.Location = &location
	}
	startsAtStr := c.FormValue("starts_at")
	if startsAtStr != "" {
		startsAt, err := time.Parse("2006-01-02 15:04:05", startsAtStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid date format",
				Details: err.Error(),
			})
		}
		req.StartsAt = &startsAt
	}

	totalTicketsStr := c.FormValue("total_tickets")
	if totalTicketsStr != "" {
		totalTickets, err := strconv.Atoi(totalTicketsStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid total tickets value",
				Details: err.Error(),
			})
		}
		req.TotalTickets = &totalTickets
	}

	priceStr := c.FormValue("price")
	if priceStr != "" {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Code:    http.StatusBadRequest,
				Message: "Invalid price value",
				Details: err.Error(),
			})
		}
		req.Price = &price
	}

	fileHeader, err := c.FormFile("photo")
	if err == nil {
		req.Photo = fileHeader
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.UpdateEvent(uint(eventId), userId, &req)

	if err != nil {
		return eventErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, response)
}
