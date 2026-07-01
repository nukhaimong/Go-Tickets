package bookings

import (
	"gotickets/internal/bookings/dto"
	"gotickets/internal/event"

	"github.com/google/uuid"
)

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository, eventRepo event.Repository) *service {
	return &service{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
}

func (s *service) CreateBooking(userId uint, req dto.CreateRequest) (*dto.Response, error) {
	booking, err := s.bookingRepo.CreateWithTicketsUpdate(userId, req.EventID, req.Quantity)
	if err != nil {
		return nil, err
	}
	return booking.ToResponse(), nil
}

func (s *service) GetMyBookings(userId uint) ([]*dto.Response, error) {
	bookings, err := s.bookingRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}
	response := make([]*dto.Response, len(bookings))

	for i, b := range bookings {
		response[i] = b.ToResponse()
	}

	return response, nil
}
