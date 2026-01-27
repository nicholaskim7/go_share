package storage

import (
	"sync"
	"time"

	"github.com/nicholaskim7/go_share/internal/models"
)

type UserStore struct {
	users []models.User
	mu    sync.Mutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: []models.User{
			{ID: 1, FirstName: "nicholas", LastName: "kim", UserName: "nkim7", Email: "nick@gmail.com", Password: "12345", DateCreated: time.Now().UTC()},
			{ID: 2, FirstName: "john", LastName: "doe", UserName: "jdoe", Email: "johndoe@gmail.com", Password: "678910", DateCreated: time.Now().UTC()},
		},
	}
}

func (s *UserStore) GetAll() []models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	// return copy to prevent accidental modifications or read while other is writing
	cp := make([]models.User, len(s.users))
	copy(cp, s.users)
	return cp
}

func (s *UserStore) Create(user models.User) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.DateCreated = time.Now().UTC()
	user.ID = int64(len(s.users) + 1)
	s.users = append(s.users, user)
	return user
}
