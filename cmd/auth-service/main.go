package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"auth-service/internal/router"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//Инициализация логгера
	logger.InitLogger()
	logger.Log.Info("Starting auth service")

	// Загрузка конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatalf("failed to load config: %v", err)
	}
	logger.Log.Info("Config loaded successfully")

	// Инициализация базы данных
	dbPool, err := db.InitDB(cfg.Postgres)
	if err != nil {
		logger.Log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbPool.Close()
	logger.Log.Info("Database connected")

	// Инициализация сервиса токенов
	refreshTokenRepo := db.NewRefreshTokenRepo(dbPool)
	tokenService := auth.NewTokenService(cfg, refreshTokenRepo)
	logger.Log.Info("Token service initialized")

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	logger.Log.Info("Auth middleware initialized")

	// Инициализация хендлера и роутера
	handler := handler.NewHandler(tokenService)
	router := router.NewRouter(handler, authMiddleware)
	logger.Log.Info("Router initialized")

	// Настройка HTTP-сервера
	addr := fmt.Sprintf("%s:%s", cfg.HTTP.ServerHost, cfg.HTTP.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ServerReadTimeout,
		WriteTimeout: cfg.HTTP.ServerWriteTimeout,
		IdleTimeout:  cfg.HTTP.ServerIdleTimeout,
	}

	// Запуск сервера
	go func() {
		logger.Log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("server forced to shutdown: %v", err)
	}
	logger.Log.Println("Server exited gracefully")

}
