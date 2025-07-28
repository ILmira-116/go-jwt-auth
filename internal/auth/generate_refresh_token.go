package auth

import (
	"crypto/rand"
	"encoding/base64"

	"auth-service/internal/logger"
)

func (s *TokenService) GenerateRefreshToken() (string, error) {
	// Генерируем случайный байтовый массив
	raw := make([]byte, refreshTokenLenght)

	_, err := rand.Read(raw)
	if err != nil {
		logger.Log.Errorf("Failed to generate refresh token: %v", err)
		return "", err
	}

	//Кодируем его в base64 (для отправки клиенту)
	token := base64.URLEncoding.EncodeToString(raw)

	// Возвращем токен
	return token, nil
}
