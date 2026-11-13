package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Password struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	GroupID   *int64    `json:"group_id,omitempty"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	URL       string    `json:"url"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordInput struct {
	GroupID  *int64 `json:"group_id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Username  string    `json:"username"`
}

type Group struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	PasswordCount int       `json:"password_count"`
}

type GroupInput struct {
	Name string `json:"name"`
}

type APIKey struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

type APIKeyInput struct {
	Name string `json:"name"`
}

type APIKeyCreated struct {
	APIKey
	Key string `json:"key"`
}
