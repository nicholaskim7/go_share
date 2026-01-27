package models

// structure of data
import "time"

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Tags        []string  `json:"tags"`
	Files       []string  `json:"files"`
	DateCreated time.Time `json:"date_created"`
}
