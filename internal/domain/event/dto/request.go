package dto

import (
	"mime/multipart"
	"time"
)

type CreateRequest struct {
	Title        string                `form:"title" validate:"required,min=2,max=150"`
	Description  string                `form:"description" validate:"required,max=1000"`
	Location     string                `form:"location" validate:"required"`
	StartsAt     time.Time             `form:"starts_at" validate:"required"`
	TotalTickets int                   `form:"total_tickets" validate:"required,gt=0"`
	Price        int                   `form:"price" validate:"gte=0"`
	Photo        *multipart.FileHeader `form:"photo" validate:"required"`
}

type UpdateRequest struct {
	Title        *string               `form:"title" validate:"omitempty,min=2,max=150"`
	Description  *string               `form:"description" validate:"omitempty,max=1000"`
	Location     *string               `form:"location" validate:"omitempty"`
	StartsAt     *time.Time            `form:"starts_at" validate:"omitempty"`
	Price        *int                  `form:"price" validate:"omitempty,gte=0"`
	TotalTickets *int                  `form:"total_tickets" validate:"omitempty,gt=0"`
	Photo        *multipart.FileHeader `form:"photo" validate:"omitempty"`
}
