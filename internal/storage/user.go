package storage

import (
	"context"

	"github.com/nicholaskim7/go_share/internal/models"
)

// describes what a user store can do, not how
type UserStore interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
}
