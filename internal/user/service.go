package user

import "gotickets/internal/user/dto"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
		// Password: req.Password,
	}

	// hash password and set to user.Password
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}
