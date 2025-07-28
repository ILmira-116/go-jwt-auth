package auth

import (
	"auth-service/internal/config"
	"auth-service/internal/db"

	"time"
)

const refreshTokenLenght = 32

type TokenService struct {
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	storage         db.RefreshTokenRepository
	webhookURL      string
}

func NewTokenService(cfg *config.Config, storage db.RefreshTokenRepository) *TokenService {
	return &TokenService{
		jwtSecret:       []byte(cfg.JWT.Secret),
		accessTokenTTL:  time.Duration(cfg.JWT.AccessTokenTTL) * time.Minute,
		refreshTokenTTL: time.Duration(cfg.JWT.RefreshTokenTTL) * time.Hour,
		storage:         storage,
		webhookURL:      cfg.WebhookURL,
	}
}
