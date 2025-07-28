package auth

import (
	"errors"
	"time"

	"auth-service/internal/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *TokenService) ValidateRefreshToken(userID uuid.UUID, token string) (bool, error) {
	records, err := s.storage.GetRefreshTokensByUserID(userID)
	if err != nil {
		logger.Log.Error("failed to get refresh tokens from DB: ", err)
		return false, err
	}

	for _, record := range records {
		err := bcrypt.CompareHashAndPassword([]byte(record.TokenHash), []byte(token))
		if err == nil {
			if record.Revoked || record.Used || time.Now().After(record.ExpiresAt) {
				logger.Log.Warn("refresh token is revoked, used or expired")
				return false, errors.New("refresh token is no longer valid")
			}
			return true, nil
		}
	}
	logger.Log.Error("refresh token hash mismatch for all tokens")
	return false, errors.New("invalid refresh token")
}
