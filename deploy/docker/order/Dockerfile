# ============================
# Stage 1: Build stage
# ============================

# Используем официальный образ Go на базе Alpine — лёгкий, быстрый и безопасный
FROM golang:1.24-alpine AS builder

# Устанавливаем git — нужен для загрузки зависимостей из приватных и публичных репозиториев
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.work и его контрольную сумму — это позволяет Go правильно связать модули
COPY go.work go.work.sum ./

# Копируем go.mod и go.sum всех зависимых модулей — нужно для кэширования go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/
COPY deploy/compose/order/.env /app/deploy/compose/order/.env

# Загружаем все зависимости, указанные в модульных файлах
RUN go mod download

# Копируем исходный код всех модулей — включая зависимости, т.к. они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY order ./order
COPY inventory ./inventory
COPY payment ./payment

COPY order/migrations ./order/migrations

# Собираем бинарный файл Inventory-сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-order ./order/cmd/order_service/main.go

# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.21.3

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Запускаем приложение под non-root пользователем

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/app-order .

# Копируем .env файл в рабочую директорию
COPY --from=builder /app/deploy/compose/order/.env ./deploy/compose/order/.env

# в DI ищет от папки корневой т.е. при обычном запуске order -> migrations а в тестах app -> migrations
COPY --from=builder /app/order/migrations ./migrations

# COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

RUN chown -R appuser:appgroup /app

USER appuser

# Экспонируем порт gRPC-сервиса Inventory
EXPOSE 8080

# Устанавливаем команду запуска — запускаем наш бинарь
ENTRYPOINT ["./app-order"]