package storage

import (
	"context"

	"github.com/nicholaskim7/go_share/internal/models"
)

// declare what a post store can do, not how it does it
type PostStore interface {
	GetAll(ctx context.Context) ([]models.Post, error)
	Create(ctx context.Context, post models.Post) (models.Post, error)
}
