package storage

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nicholaskim7/go_share/internal/models"
)

type PostDBStore struct {
	db *sql.DB
}

func NewPostDBStore(db *sql.DB) *PostDBStore {
	return &PostDBStore{db: db}
}

func (s *PostDBStore) GetAll(ctx context.Context) ([]models.Post, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, body, tags, files, date_created FROM posts`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Body,
			pq.Array(&p.Tags),
			pq.Array(&p.Files),
			&p.DateCreated,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *PostDBStore) Create(ctx context.Context, post models.Post) (models.Post, error) {
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO posts (user_id, title, body, tags, files)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, date_created`,
		post.UserID,
		post.Title,
		post.Body,
		pq.Array(post.Tags),
		pq.Array(post.Files),
	).Scan(&post.ID, &post.DateCreated)
	if err != nil {
		return models.Post{}, err
	}
	return post, err
}

func (s *PostDBStore) GetByUsername(ctx context.Context, username string) ([]models.Post, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT p.id, p.user_id, p.title, p.body, p.tags, p.files, p.date_created 
		 FROM posts p
		 JOIN users u ON p.user_id = u.id
		 WHERE u.user_name = $1`,
		username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Body,
			pq.Array(&p.Tags),
			pq.Array(&p.Files),
			&p.DateCreated,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
