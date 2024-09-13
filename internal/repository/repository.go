package repository

import (
	"context"

	"shaffra_assessment/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, client *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUserByID(ctx context.Context, id string, user *models.User) error
	DeleteUserID(ctx context.Context, id string) error
}
