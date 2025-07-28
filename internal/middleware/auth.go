package middleware

import (
	"auth-service/internal/auth"
	"auth-service/internal/logger"
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type AuthMiddleware struct {
	tokenService *auth.TokenService
}

func NewAuthMiddleware(ts *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: ts}
}

// Middleware проверяет Access токен в заголовке Authorization
func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Warn("AuthMiddleware: missing Authorization header")
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		// Формат заголовка должен быть: "Bearer {токен}"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Log.Warnf("AuthMiddleware: invalid Authorization header format: %s", authHeader)
			http.Error(w, "authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1] // вытаскиваем сам токен

		// Проверяем токен с помощью сервиса
		userID, err := a.tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			logger.Log.Warnf("AuthMiddleware: failed to validate access token: %v", err)
			http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
			return
		}

		// Если токен валидный — сохраняем userID в контекст запроса,
		// чтобы дальше в хендлерах можно было получить ID пользователя
		ctx := r.Context()
		ctx = contextWithUserID(ctx, userID)

		// Передаём запрос дальше, но с обновлённым контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Ключ для хранения userID в контексте
type ctxKeyUserID struct{}

// contextWithUserID сохраняет userID в контекст
func contextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

// UserIDFromContext достаёт userID из контекста
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(ctxKeyUserID{}).(uuid.UUID)
	return userID, ok
}
