package models

// structure of data
import "time"

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateCreated time.Time `json:"date_created"`
}

type UserLoginPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
