package bookings

import (
	"fmt"
	"gotickets/internal/domain/bookings/dto"
	"gotickets/internal/domain/event"
	"gotickets/internal/payment"
	"strconv"

	"github.com/google/uuid"
)

var (
	bookingPending   = "pending"
	bookingConfirmed = "confirmed"
)

type service struct {
	bookingRepo   Repository
	eventRepo     event.Repository
	stripeService *payment.StripeService
}

func NewService(bookingRepo Repository, eventRepo event.Repository, stripeService *payment.StripeService) *service {
	return &service{
		bookingRepo:   bookingRepo,
		eventRepo:     eventRepo,
		stripeService: stripeService,
	}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
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

func(s *service) GetByID(bookingId uint) (*dto.Response, error) {
	booking, err := s.bookingRepo.GetByID(bookingId)
	if err != nil {
		return nil, err
	}
	return booking.ToResponse(), nil
}

func (s *service) CreateBooking(userId uint, req dto.CreateRequest) (*dto.Response, error) {
	// 1. Get event
	eventData, err := s.eventRepo.GetEventByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if eventData.AvailableTickets < req.Quantity {
		return nil, ErrNotEnoughTickets
	}

	// 2. Create booking with "pending" status
	booking := &Booking{
		UserID:      userId,
		EventID:     req.EventID,
		Quantity:    req.Quantity,
		TotalPrice:  eventData.Price * req.Quantity,
		Status:      bookingPending,
		BookingCode: generateBookingCode(),
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}

	// 3. Create Stripe checkout session
	checkoutReq := &payment.CreateCheckoutRequest{
		BookingID: booking.ID,
		UserID:    userId,
		EventID:   req.EventID,
		Amount:    int64(eventData.Price),
		Quantity:  req.Quantity,
	}

	checkoutResp, err := s.stripeService.CreateCheckoutSession(checkoutReq)
	if err != nil {
		s.bookingRepo.DeleteBooking(booking.ID)
		return nil, err
	}

	// 4. Save session ID to booking
	booking.StripeSessionID = checkoutResp.SessionID
	s.bookingRepo.Update(booking)

	// 5. Return the checkout URL to the client
	response := booking.ToResponse()
	response.CheckoutURL = checkoutResp.URL

	return response, nil
}

// HandlePaymentSuccess handles successful payment webhook
func (s *service) HandlePaymentSuccess(bookingID string) error {
	id, err := strconv.ParseUint(bookingID, 10, 32)
	if err != nil {
		return err
	}

	booking, err := s.bookingRepo.GetByID(uint(id))
	if err != nil {
		return err
	}

	booking.Status = bookingConfirmed

	eventData, err := s.eventRepo.GetEventByID(booking.EventID)
	if err != nil {
		return err
	}

	eventData.AvailableTickets -= booking.Quantity
	if err := s.eventRepo.Update(eventData); err != nil {
		return err
	}

	return s.bookingRepo.Update(booking)
}

// HandlePaymentExpired handles expired payment webhook
func (s *service) HandlePaymentExpired(bookingID string) error {
	// Convert string to uint
	id, err := strconv.ParseUint(bookingID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid booking ID: %w", err)
	}

	// Delete the booking
	return s.bookingRepo.DeleteBooking(uint(id))
}
