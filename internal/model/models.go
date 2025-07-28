package model

import (
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

type RefreshTokenRecord struct {
	UserID    uuid.UUID
	TokenHash string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Revoked   bool
	Used      bool
	UserAgent string
	IP        string
}
