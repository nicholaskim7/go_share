package storage

import (
	"context"
	"sync"
	"time"

	"github.com/nicholaskim7/go_share/internal/models"
)

type UserMemoryStore struct {
	users []models.User
	mu    sync.Mutex
}

func NewUserMemoryStore() *UserMemoryStore {
	return &UserMemoryStore{
		users: []models.User{
			{ID: 1, FirstName: "nicholas", LastName: "kim", UserName: "nkim7", Email: "nick@gmail.com", Password: "12345", DateCreated: time.Now().UTC()},
			{ID: 2, FirstName: "john", LastName: "doe", UserName: "jdoe", Email: "johndoe@gmail.com", Password: "678910", DateCreated: time.Now().UTC()},
		},
	}
}

func (s *UserMemoryStore) GetAll(_ context.Context) ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// return copy to prevent accidental modifications or read while other is writing
	cp := make([]models.User, len(s.users))
	copy(cp, s.users)
	return cp, nil
}

func (s *UserMemoryStore) Create(_ context.Context, user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.DateCreated = time.Now().UTC()
	user.ID = int64(len(s.users) + 1)
	s.users = append(s.users, user)
	return user, nil
}
