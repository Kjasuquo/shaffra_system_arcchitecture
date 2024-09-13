package service

import (
	"context"
	"fmt"
	"strings"

	"shaffra_assessment/internal/models"
)

// CreateUser has the logic for creating a user
func (s *service) CreateUser(ctx context.Context, user *models.User) (string, error) {
	id, err := s.user.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", fmt.Errorf("service---email already exists")
		}
		return "", fmt.Errorf("service---create user: %v", err)
	}

	return id, nil
}

// GetUserByID has the logic for retrieving a user by ID
func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.user.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service---get user: %v", err)
	}

	return user, nil
}

// UpdateUserByID has the logic for updating a user by ID
func (s *service) UpdateUserByID(ctx context.Context, id string, user *models.User) error {

	u, err := s.user.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("service---get user: %v", err)
	}
	if user.Name != "" {
		u.Name = user.Name
	}
	if user.Email != "" {
		u.Email = user.Email
	}
	if user.Age != 0 {
		u.Age = user.Age
	}

	err = s.user.UpdateUserByID(ctx, id, u)
	if err != nil {
		return fmt.Errorf("service---update user: %v", err)
	}

	return nil
}

// DeleteUserID has the logic for deleting a user by ID
func (s *service) DeleteUserID(ctx context.Context, id string) error {
	err := s.user.DeleteUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("service---delete user: %v", err)
	}

	return nil
}
