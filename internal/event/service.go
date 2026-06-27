package event

import "gotickets/internal/event/dto"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateEvent(req dto.CreateRequest) (*dto.Response, error) {
	event := &Event{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
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
