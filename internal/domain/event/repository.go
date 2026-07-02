package event

import (
	"errors"

	"gorm.io/gorm"
)

var ErrEventNotFound = errors.New("Event not found")

type Repository interface {
	Create(event *Event) error
	GetAll() ([]*Event, error)
	GetEventByID(eventId uint) (*Event, error)
	Update(event *Event) error
}

type respository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &respository{
		db: db,
	}
}

func (r *respository) Create(event *Event) error {
	return r.db.Create(event).Error
}

func (r *respository) GetAll() ([]*Event, error) {
	var events []*Event

	err := r.db.Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *respository) GetEventByID(eventId uint) (*Event, error) {
	event := &Event{}

	err := r.db.First(event, eventId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (r *respository) Update(event *Event) error {
	return r.db.Save(event).Error
}
