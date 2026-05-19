package user

import (
	"fmt"

	"backend-travelku/pkg/auth"
	"backend-travelku/pkg/logger"
)

// Service adalah core business logic untuk fitur user.
type Service interface {
	GetUserByID(id string) (*User, error)
	GetAllUsers() ([]User, error)
	UpdateProfile(id string, req *UpdateProfileRequest) (*User, error)
	ChangePassword(id string, req *ChangePasswordRequest) error
	DeleteUser(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserByID(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		logger.Warn("User not found: " + id)
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *service) GetAllUsers() ([]User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		logger.Error("Failed to get users: " + err.Error())
		return nil, fmt.Errorf("failed to get users")
	}
	return users, nil
}

func (s *service) UpdateProfile(id string, req *UpdateProfileRequest) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		logger.Warn("User not found: " + id)
		return nil, fmt.Errorf("user not found")
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if user.Email != req.Email {
		if s.repo.ExistsByEmail(req.Email) {
			logger.Warn("Email already registered: " + req.Email)
			return nil, fmt.Errorf("email already registered")
		}
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		logger.Error("Failed to update profile: " + err.Error())
		return nil, fmt.Errorf("failed to update profile")
	}

	logger.Info("User profile updated successfully: " + id)
	return user, nil
}

func (s *service) ChangePassword(id string, req *ChangePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		logger.Warn("User not found: " + id)
		return fmt.Errorf("user not found")
	}

	if !auth.VerifyPassword(user.Password, req.OldPassword) {
		logger.Warn("Change password failed: invalid old password - " + id)
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("Failed to hash password")
		return fmt.Errorf("failed to hash password")
	}

	user.Password = hashedPassword
	if err := s.repo.Update(user); err != nil {
		logger.Error("Failed to change password: " + err.Error())
		return fmt.Errorf("failed to change password")
	}

	logger.Info("User password changed successfully: " + id)
	return nil
}

func (s *service) DeleteUser(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		logger.Warn("User not found: " + id)
		return fmt.Errorf("user not found")
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Error("Failed to delete user: " + err.Error())
		return fmt.Errorf("failed to delete user")
	}

	logger.Info("User deleted successfully: " + id)
	return nil
}
