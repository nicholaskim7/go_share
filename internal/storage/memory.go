package storage

// db functions added in the future
import (
	"sync"
	"time"

	"github.com/nicholaskim7/go_share/internal/models"
)

type PostStore struct {
	posts []models.Post
	mu    sync.Mutex
}

func NewPostStore() *PostStore {
	return &PostStore{
		posts: []models.Post{
			{ID: 1, UserID: 1, Title: "python script", Body: "code contains script", Tags: []string{"programming", "python", "script"}, Files: []string{"script.py"}, DateCreated: time.Now().UTC()},
			{ID: 2, UserID: 1, Title: "go script", Body: "code contains calculator app", Tags: []string{"coding", "python", "go"}, Files: []string{"main.go"}, DateCreated: time.Now().UTC()},
		},
	}
}

func (s *PostStore) GetAll() []models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	// return copy to prevent accidental modifications or read while other is writing
	cp := make([]models.Post, len(s.posts))
	copy(cp, s.posts)
	return cp
}

func (s *PostStore) Create(post models.Post) models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	post.DateCreated = time.Now().UTC()
	post.ID = int64(len(s.posts) + 1)
	s.posts = append(s.posts, post)
	return post
}
