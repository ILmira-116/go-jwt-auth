package auth

import (
	"auth-service/internal/logger"

	"github.com/google/uuid"
)

func (s *TokenService) RevokeAllTokensForUser(userID uuid.UUID) error {
	// Обновить в БД все токены пользователя, выставив Revoked = true
	err := s.storage.RevokeTokensByUserID(userID)
	if err != nil {
		logger.Log.Error("failed to revoke tokens for user: ", err)
		return err
	}
	return nil
}

func (s *TokenService) RevokeUser(userID uuid.UUID) error {
	return s.storage.RevokeUser(userID)
}
