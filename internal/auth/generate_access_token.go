package auth

import (
	"time"

	"auth-service/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *TokenService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	// Формируем payload
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(s.accessTokenTTL).Unix(), // срок жизни токена
		"iat": time.Now().Unix(),                       // время создания токена
	}

	// Создаём новый JWT токен с указанными claims и алгоритмом подписи HS512 (HMAC с SHA-512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Подписываем токен
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		logger.Log.Errorf("Failed to sign access token for user %s: %v", userID, err)
		return "", err
	}

	// Возвращаем готовый подписанный JWT токен в виде строки
	logger.Log.Infof("Access token successfully generated for user %s", userID)
	return signedToken, nil
}
