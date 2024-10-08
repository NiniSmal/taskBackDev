package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

var ErrNotFound = errors.New("not found")
var ErrUserBlocked = errors.New("user blocked")

type Session struct {
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

var JwtSecretKey = "your-256-bit-secret"
