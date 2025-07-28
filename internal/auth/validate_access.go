package auth

import (
	"errors"

	"auth-service/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *TokenService) ValidateAccessToken(tokenStr string) (uuid.UUID, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		logger.Log.Error("access token is invalid: ", err)
		return uuid.Nil, errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Log.Error("invalid claims structure in token")
		return uuid.Nil, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		logger.Log.Error("token does not contain subject (user_id)")
		return uuid.Nil, errors.New("no subject in token")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		logger.Log.Error("failed to parse user_id from token subject: ", err)
		return uuid.Nil, errors.New("invalid user_id format in token")
	}

	revoked, err := s.storage.IsUserRevoked(userID)
	if err != nil {
		logger.Log.Error("failed to check revoked users: ", err)
		return uuid.Nil, errors.New("internal error")
	}
	if revoked {
		return uuid.Nil, errors.New("user is revoked")
	}

	return userID, nil
}
