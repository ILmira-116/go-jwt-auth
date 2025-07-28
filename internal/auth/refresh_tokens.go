package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"auth-service/internal/logger"
	"auth-service/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *TokenService) RefreshTokens(userID uuid.UUID, refreshToken, userAgent, ip string) (model.Tokens, error) {
	records, err := s.storage.GetRefreshTokensByUserID(userID)
	if err != nil {
		logger.Log.Error("failed to get refresh tokens from DB: ", err)
		return model.Tokens{}, err
	}

	var validRecord *model.RefreshTokenRecord
	for i := range records {
		record := &records[i]
		// Проверяем совпадение хэша
		if err := bcrypt.CompareHashAndPassword([]byte(record.TokenHash), []byte(refreshToken)); err == nil {
			// Проверяем статус токена
			if record.Revoked || record.Used || time.Now().After(record.ExpiresAt) {
				continue
			}
			validRecord = record
			break
		}
	}

	if validRecord == nil {
		logger.Log.Warn("refresh token validation failed")
		return model.Tokens{}, errors.New("invalid or expired refresh token")
	}

	// Проверяем совпадение User-Agent
	if validRecord.UserAgent != userAgent {
		err := s.RevokeAllTokensForUser(userID)
		if err != nil {
			logger.Log.Error("failed to revoke all tokens: ", err)
		}
		return model.Tokens{}, errors.New("user-agent mismatch — user deauthorized")
	}

	// Проверяем IP, если он изменился — вызываем webhook
	if validRecord.IP != ip {
		go func() {
			payload := map[string]string{
				"user_id": userID.String(),
				"old_ip":  validRecord.IP,
				"new_ip":  ip,
				"time":    time.Now().Format(time.RFC3339),
			}
			body, _ := json.Marshal(payload)
			http.Post(s.webhookURL, "application/json", bytes.NewReader(body))
		}()
	}

	// Отмечаем токен как использованный
	validRecord.Used = true
	if err := s.storage.UpdateRefreshToken(userID, *validRecord); err != nil {
		logger.Log.Error("failed to update refresh token: ", err)
		return model.Tokens{}, err
	}

	// Генерируем новую пару токенов
	tokens, err := s.GenerateTokenPair(userID, userAgent, ip)
	if err != nil {
		logger.Log.Error("failed to generate new tokens: ", err)
		return model.Tokens{}, err
	}

	return *tokens, nil
}
