package bookings

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrBookingNotFound         = errors.New("Booking not found")
	ErrNotEnoughTickets        = errors.New("Not enough tickets available")
	ErrBookingAlreadyCancelled = errors.New("Booking already cancelled")
	ErrForbiddenBookingAccess  = errors.New("You do not own this booking")
)

type Repository interface {
	Create(booking *Booking) error
	GetByID(bookingId uint) (*Booking, error)
	GetByUserID(userId uint) ([]*Booking, error)
	Update(booking *Booking) error
	//CreateWithTicketsUpdate(userId uint, eventId uint, quantity int) (*Booking, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepostory(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(booking *Booking) error {
	return r.db.Create(booking).Error
}

func (r *repository) GetByID(bookingId uint) (*Booking, error) {
	booking := &Booking{}

	err := r.db.First(booking, bookingId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return booking, nil
}

func (r *repository) GetByUserID(userId uint) ([]*Booking, error) {
	var bookings []*Booking

	err := r.db.Where("user_id= ?", userId).Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}
