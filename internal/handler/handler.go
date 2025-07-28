package handler

import (
	"auth-service/internal/auth"
)

type Handler struct {
	tokenService *auth.TokenService
}

func NewHandler(tokenService *auth.TokenService) *Handler {
	return &Handler{tokenService: tokenService}
}
