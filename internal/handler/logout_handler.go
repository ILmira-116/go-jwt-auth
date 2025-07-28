package handler

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"auth-service/internal/pkg"
)

// LogoutHandler godoc
// @Summary Деавторизация пользователя (logout)
// @Description Отзывает все refresh токены пользователя и запрещает дальнейший доступ к защищенным эндпоинтам
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} pkg.MessageResponse "Сообщение об успешном выходе"
// @Failure 401 {object} model.ErrorResponse "user not authorized"
// @Failure 500 {object} model.ErrorResponse "failed to revoke tokens или failed to revoke user access"
// @Router /logout [post]
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		logger.Log.Warn("LogoutHandler: userID not found in context")
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	// Отзываем все refresh токены
	if err := h.tokenService.RevokeAllTokensForUser(userID); err != nil {
		logger.Log.Errorf("LogoutHandler: failed to revoke tokens for user %s: %v", userID, err)
		http.Error(w, "failed to revoke tokens", http.StatusInternalServerError)
		return
	}

	// Добавляем пользователя в revoked_users
	if err := h.tokenService.RevokeUser(userID); err != nil {
		logger.Log.Errorf("LogoutHandler: failed to revoke user %s: %v", userID, err)
		http.Error(w, "failed to revoke user access", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	resp := pkg.MessageResponse{Message: "user logged out successfully"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Errorf("LogoutHandler: failed to encode response for user %s: %v", userID, err)
		pkg.WriteJSONError(w, "failed to respond", http.StatusInternalServerError)
	}
}
