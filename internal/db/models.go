package db

import "time"

type User struct {
	
	PublicID string    `json:"public_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Omit password in JSON responses
}

type Post struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	PosterID int       `json:"poster"`
	PostedAt time.Time `json:"posted_at"`
}