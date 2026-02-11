package storage

import (
	"context"

	"github.com/nicholaskim7/go_share/internal/models"
)

// declare what a post store can do, not how it does it
type PostStore interface {
	GetAll(ctx context.Context) ([]models.Post, error)
	Create(ctx context.Context, post models.Post) (models.Post, error)
	Delete(ctx context.Context, postId int64, userId int64) error
	GetByUsername(ctx context.Context, username string) ([]models.Post, error)
	GetByTag(ctx context.Context, tag string) ([]models.Post, error)
}
