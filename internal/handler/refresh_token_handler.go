package handler

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"auth-service/internal/pkg"
)

// RefreshTokenHandler godoc
// @Summary      Обновление пары токенов (access и refresh)
// @Description  Обновляет access и refresh токены по действительному refresh токену
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        refresh_token body pkg.RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} model.Tokens "Новая пара токенов"
// @Failure      400 {object} model.ErrorResponse "invalid request body"
// @Failure      401 {object} model.ErrorResponse "user not authorized или failed to refresh tokens"
// @Router       /token/refresh [post]
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		logger.Log.Warn("RefreshTokenHandler: userID not found in context")
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	// Декодируем тело запроса
	var req pkg.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Warnf("RefreshTokenHandler: failed to decode body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	logger.Log.Infof("RefreshTokenHandler: userID=%s, ip=%s, agent=%s, refresh_token=%s",
		userID.String(), ip, userAgent, req.RefreshToken)

	// Пытаемся обновить токены
	tokens, err := h.tokenService.RefreshTokens(userID, req.RefreshToken, userAgent, ip)
	if err != nil {
		logger.Log.Errorf("RefreshTokenHandler: failed to refresh tokens: %v", err)
		http.Error(w, "failed to refresh tokens", http.StatusUnauthorized)
		return
	}

	logger.Log.Infof("RefreshTokenHandler: successfully refreshed tokens for userID=%s", userID.String())

	// Отправляем обновлённые токены
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logger.Log.Errorf("RefreshTokenHandler: failed to encode tokens: %v", err)
		pkg.WriteJSONError(w, "internal error", http.StatusInternalServerError)
	}
}
