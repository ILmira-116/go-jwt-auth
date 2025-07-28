package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"net/http"

	_ "auth-service/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.Handler, authMW *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Post("/token", h.GenerateTokenHandler)

	r.Group(func(r chi.Router) {
		r.Use(authMW.Middleware)
		r.Post("/token/refresh", h.RefreshTokenHandler)
		r.Get("/users/me", h.GetCurrentUserIDHandler)
		r.Post("/logout", h.LogoutHandler)
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}
