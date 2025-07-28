package handler

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/pkg"

	"github.com/google/uuid"
)

// GenerateTokenHandler godoc
// @Summary      Получить access и refresh токены
// @Description  Генерация токенов по переданному user_id (UUID)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_id query string true "UUID пользователя"
// @Success      200 {object} model.Tokens "Пара access и refresh токенов"
// @Failure      400 {object} model.ErrorResponse "Неверный запрос"
// @Failure      500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /token [post]
func (h *Handler) GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		pkg.WriteJSONError(w, "user_id is required", http.StatusBadRequest)
		logger.Log.Warn("GenerateTokenHandler: missing user_id in query")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		pkg.WriteJSONError(w, "invalid user_id", http.StatusBadRequest)
		logger.Log.Warnf("GenerateTokenHandler: invalid user_id format: %v", err)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	tokens, err := h.tokenService.GenerateTokenPair(userID, userAgent, ip)
	if err != nil {
		pkg.WriteJSONError(w, "failed to generate tokens", http.StatusInternalServerError)
		logger.Log.Errorf("GenerateTokenHandler: failed to generate token pair: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logger.Log.Errorf("GenerateTokenHandler: failed to encode response: %v", err)
		pkg.WriteJSONError(w, "internal error", http.StatusInternalServerError)
	}
}
