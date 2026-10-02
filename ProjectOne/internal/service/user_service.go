package service

import (
	"appMove/internal/dto"
	"appMove/internal/model"
	"appMove/internal/repository"
	"appMove/pkg/auth"
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {

	//Основные
	Register(ctx context.Context, reg *dto.RequestRegister) (*dto.UserProfileResponse, error)
	Login(ctx context.Context, log *dto.RequestLogin) (string, error)
	Logout(ctx context.Context) error
	GetProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error)

	//Изменения
	ChangePassword(ctx context.Context) error

	//Дополнительно
}

type userService struct {
	repo         repository.UserRepository
	tokenManager auth.TokenManager
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error) {
	const op = "Service.GetProfile"

	if userID <= 0 {
		return nil, fmt.Errorf("%s: invalid userID %d", op, userID)
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return &dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) Register(ctx context.Context, reg *dto.RequestRegister) (*dto.UserProfileResponse, error) {
	const op = "Service.Register"

	if reg.Email == "" {
		return nil, fmt.Errorf("%s: email is required: %w", op, ErrInvalidInput)
	}

	if utf8.RuneCountInString(reg.Password) < 8 {
		return nil, fmt.Errorf("%s: password must be at least 8 characters: %w", op, ErrInvalidInput)

	}

	if reg.Name == "" {
		return nil, fmt.Errorf("%s: name is required: %w", op, ErrInvalidInput)
	}

	userEmail, err := s.repo.GetByEmail(ctx, reg.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if userEmail != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrEmailAlreadyExists)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userModel := &model.User{
		Email:        reg.Email,
		PasswordHash: string(passHash),
		Name:         reg.Name,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userReg, err := s.repo.Create(ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.UserProfileResponse{
		ID:        userReg,
		Email:     userModel.Email,
		Name:      userModel.Name,
		Role:      userModel.Role,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}, nil
}

func (s *userService) Login(ctx context.Context, log *dto.RequestLogin) (string, error) {
	const op = "Service.Login"

	if log.Email == "" {
		return "", fmt.Errorf("%s: email is required: %w", op, ErrInvalidInput)
	}

	if log.Password == "" {
		return "", fmt.Errorf("%s: password is required: %w", op, ErrInvalidInput)
	}

	user, err := s.repo.GetByEmail(ctx, log.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(log.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: internal bcrypt error: %w", op, err)
	}

	token, err := s.tokenManager.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *userService) ChangePassword(ctx context.Context) error {
	const op = "Service.ChangePassword"

	return nil
}

func (s *userService) Logout(ctx context.Context) error {
	const op = "Service.Logout"

	return nil
}
