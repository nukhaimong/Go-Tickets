package user

import (
	"fmt"
	"gotickets/internal/auth"
	"gotickets/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("Invalid Email or Password")

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
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

	// generate toke
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = user.checkPassword(req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// generate toke
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token: %w", err)
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}
	return &response, nil
}
