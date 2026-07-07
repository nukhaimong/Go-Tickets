package event

import (
	"errors"
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/event/dto"
)

type service struct {
	repo       Repository
	cloudinary *config.CloudinaryService
}

func NewService(repo Repository, cloudinary *config.CloudinaryService) *service {
	return &service{
		repo:       repo,
		cloudinary: cloudinary,
	}
}

func (s *service) CreateEvent(req *dto.CreateRequest, userId uint) (*dto.Response, error) {
	event := &Event{
		UserID:           userId,
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
		PhotoURL:         "",
	}

	// upload the photo if provided
	if req.Photo != nil {
		file, err := req.Photo.Open()
		if err != nil {
			return nil, errors.New("Failed to open file: " + err.Error())
		}
		defer file.Close()
		// upload file to cloudinary
		photoURL, err := s.cloudinary.UploadEventImage(file, req.Photo)
		if err != nil {
			return nil, errors.New("failed to upload photo: " + err.Error())
		}
		event.PhotoURL = photoURL
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	return event.ToResponse(), nil
}

func (s *service) GetEvents() ([]*dto.Response, error) {

	events, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	responses := make([]*dto.Response, len(events))
	for i, event := range events {
		responses[i] = event.ToResponse()
	}
	return responses, nil
}

func (s *service) GetEventById(eventId uint) (*dto.Response, error) {

	event, err := s.repo.GetEventByID(eventId)
	if err != nil {
		return nil, err
	}
	return event.ToResponse(), nil
}

func (s *service) UpdateEvent(eventId uint, userId uint, req *dto.UpdateRequest) (*dto.Response, error) {
	event, err := s.repo.GetEventByID(eventId)
	if err != nil {
		return nil, err
	}

	if userId != event.UserID {
		return nil, errors.New("Uauthorize to update the event")
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.Location != nil {
		event.Location = *req.Location
	}

	if req.Price != nil {
		event.Price = *req.Price
	}
	if req.StartsAt != nil && !req.StartsAt.IsZero() {
		event.StartsAt = *req.StartsAt
	}
	if req.TotalTickets != nil {
		// Adjust available tickets based on the change in total tickets
		difference := *req.TotalTickets - event.TotalTickets
		event.AvailableTickets += difference
		event.TotalTickets = *req.TotalTickets
	}

	if req.Photo != nil {
		file, err := req.Photo.Open()
		if err != nil {
			return nil, errors.New("Failed to open file: " + err.Error())
		}
		defer file.Close()

		// delete the old photo from cloudinary if it exists
		if event.PhotoURL != "" {
			err := s.cloudinary.DeleteEventImage(event.PhotoURL)
			if err != nil {
				return nil, errors.New("failed to delete old photo: " + err.Error())
			}
			fmt.Println("successfully deleted photo: ", event.PhotoURL)
		}
		// upload file to cloudinary
		photoURL, err := s.cloudinary.UploadEventImage(file, req.Photo)
		if err != nil {
			return nil, errors.New("failed to upload photo: " + err.Error())
		}
		event.PhotoURL = photoURL
	}

	if err := s.repo.Update(event); err != nil {
		return nil, err
	}

	return event.ToResponse(), nil
}
