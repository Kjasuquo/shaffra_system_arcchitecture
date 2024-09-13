package service

import (
	"context"
	"shaffra_assessment/internal/repository"

	"shaffra_assessment/internal/models"
)

type service struct {
	user repository.UserRepository
}

// NewService is a constructor for service
func NewService(db repository.UserRepository) *service {
	return &service{
		user: db,
	}
}

// UserService interface has the user's logic methods
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUserByID(ctx context.Context, id string, user *models.User) error
	DeleteUserID(ctx context.Context, id string) error
}
