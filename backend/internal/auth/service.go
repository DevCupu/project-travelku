package auth

import (
	"fmt"
	"time"

	userfeature "backend-travelku/internal/user"
	jwtauth "backend-travelku/pkg/auth"
	"backend-travelku/pkg/logger"

	"github.com/google/uuid"
)

// Service adalah core business logic untuk fitur auth.
type Service interface {
	Register(req *RegisterRequest) (*userfeature.User, error)
	Login(req *LoginRequest) (string, *userfeature.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(req *RegisterRequest) (*userfeature.User, error) {
	if s.repo.ExistsByEmail(req.Email) {
		logger.Warn("Registration failed: email already registered")
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := jwtauth.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password: " + err.Error())
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &userfeature.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		logger.Error("Failed to create user: " + err.Error())
		return nil, fmt.Errorf("failed to register user")
	}

	logger.Info("User registered successfully: " + user.Email)
	return user, nil
}

func (s *service) Login(req *LoginRequest) (string, *userfeature.User, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		logger.Warn("Login failed: user not found - " + req.Email)
		return "", nil, fmt.Errorf("invalid email or password")
	}

	if !user.IsActive {
		logger.Warn("Login failed: user inactive - " + req.Email)
		return "", nil, fmt.Errorf("user account is inactive")
	}

	if !jwtauth.VerifyPassword(user.Password, req.Password) {
		logger.Warn("Login failed: invalid password - " + req.Email)
		return "", nil, fmt.Errorf("invalid email or password")
	}

	token, err := jwtauth.GenerateToken(user.ID, user.Email, 24)
	if err != nil {
		logger.Error("Failed to generate token: " + err.Error())
		return "", nil, fmt.Errorf("failed to generate token")
	}

	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.Update(user); err != nil {
		logger.Error("Failed to update last login: " + err.Error())
	}

	logger.Info("User logged in successfully: " + user.Email)
	return token, user, nil
}
