package auth

import (
	"auth-service/internal/logger"
	"auth-service/internal/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *TokenService) GenerateTokenPair(userID uuid.UUID, userAgent, ip string) (*model.Tokens, error) {
	// 1. Генерация access токена
	accessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		logger.Log.Error("failed to generate access token: ", err)
		return nil, err
	}

	// 2. Генерация refresh токена
	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		logger.Log.Error("failed to hash refresh token: ", err)
		return nil, err
	}

	// 3. Хешируем refresh токен для безопасного хранения
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	// 4. Сохраняем хеш в БД
	record := model.RefreshTokenRecord{
		UserID:    userID,
		TokenHash: string(hashedToken),
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		Revoked:   false,
		Used:      false,
		UserAgent: userAgent,
		IP:        ip,
	}

	err = s.storage.SaveRefreshToken(record)

	if err != nil {
		logger.Log.Error("failed to store refresh token: ", err)
		return nil, err
	}

	// 5. Возвращаем пару токенов
	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
