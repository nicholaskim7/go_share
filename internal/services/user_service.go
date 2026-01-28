package services

import (
	"context"

	"github.com/nicholaskim7/go_share/internal/models"
	"github.com/nicholaskim7/go_share/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store storage.UserStore
}

func NewUserService(store storage.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.store.GetAll(ctx)
}

// hash password string and return string hash
func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashed)
	return s.store.Create(ctx, user)
}
