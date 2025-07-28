# 1. Используем официальный образ Golang для сборки
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код приложения
COPY . .

# Собираем бинарник
RUN go build -o auth-service ./cmd/auth-service

# 2. Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/auth-service .

# Копируем .env, если нужно
COPY .env .

# Указываем порт, если нужно
EXPOSE 8080

# Запускаем приложение
CMD ["./auth-service"]
