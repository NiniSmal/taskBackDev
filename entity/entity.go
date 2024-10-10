package entity

import (
	"errors"
	"time"
)

type TokenResponse struct {
	Name         string
	AccessToken  string
	RefreshToken int64
}

var ErrNotFound = errors.New("not found")
var ErrUserBlocked = errors.New("user blocked")
var ErrorValidate = errors.New("validate")

type Session struct {
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsUsed       bool      `json:"is_used"`
}

var JwtSecretKey = "your-256-bit-secret"

type Email struct {
	From    string `json:"from"`
	Text    string `json:"text"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	File    string `json:"file"`
}
