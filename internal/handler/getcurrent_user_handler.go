package handler

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"auth-service/internal/pkg"
)

// GetCurrentUserIDHandler godoc
// @Summary      Получить GUID текущего пользователя
// @Description  Возвращает GUID пользователя из контекста авторизации (access token должен быть валиден)
// @Tags         users
// @Produce      json
// @Success      200 {object} pkg.UserIDResponse "GUID текущего пользователя"
// @Failure      401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure      500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *Handler) GetCurrentUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		logger.Log.Warn("GetCurrentUserIDHandler: userID not found in context")
		pkg.WriteJSONError(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	resp := pkg.UserIDResponse{GUID: userID.String()}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Errorf("GetCurrentUserIDHandler: failed to write response: %v", err)
		pkg.WriteJSONError(w, "failed to respond", http.StatusInternalServerError)
		return
	}
}
