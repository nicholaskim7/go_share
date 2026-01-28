package storage

import (
	"context"
	"database/sql"

	"github.com/nicholaskim7/go_share/internal/models"
)

type UserDBStore struct {
	db *sql.DB
}

func NewUserDBStore(db *sql.DB) *UserDBStore {
	return &UserDBStore{db: db}
}

func (s *UserDBStore) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, first_name, last_name, user_name, email, date_created FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.UserName,
			&u.Email,
			&u.DateCreated,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// ID and date_created will be done by the database
func (s *UserDBStore) Create(ctx context.Context, user models.User) (models.User, error) {
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (first_name, last_name, user_name, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, date_created`,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.DateCreated)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserDBStore) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, first_name, last_name, user_name, email, password, date_created
		 FROM users WHERE user_name = $1`,
		username,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.DateCreated,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
