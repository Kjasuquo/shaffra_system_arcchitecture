package repository

import (
	"context"
	"fmt"
	"shaffra_assessment/internal/models"
)

// CreateUser creates user in the database
func (p *PostgresDB) CreateUser(ctx context.Context, client *models.User) (string, error) {
	err := p.instance.WithContext(ctx).Create(client).Error
	if err != nil {
		return "", fmt.Errorf("error creating user: %s", err)
	}
	return client.ID, nil
}

// GetUserByID retrieves a user from the database by ID
func (p *PostgresDB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	client := &models.User{}
	if err := p.instance.WithContext(ctx).Where("id =?", id).First(client).Error; err != nil {
		return nil, fmt.Errorf("error retrieving user: %s", err)
	}
	return client, nil
}

// UpdateUserByID updates a user in the database by ID
func (p *PostgresDB) UpdateUserByID(ctx context.Context, id string, user *models.User) error {
	if err := p.instance.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).
		Updates(user).Error; err != nil {
		return fmt.Errorf("error updating user: %s", err)
	}

	return nil
}

// DeleteUserID deletes a user in the database by ID
func (p *PostgresDB) DeleteUserID(ctx context.Context, id string) error {
	news := &models.User{}
	err := p.instance.WithContext(ctx).Delete(news, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("error deleting user: %s", err)
	}
	return nil
}
