# Используем официальный образ Go на основе Alpine
FROM golang:1.24.5-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /product_accounting_service ./cmd/app/

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /

# Копируем скомпилированное приложение из предыдущего образа
COPY --from=builder /product_accounting_service /product_accounting_service

# Копируем файлы миграции
COPY migrations /migrations
COPY .env .

# Открываем порт
EXPOSE 8080
ENV PORT=8080

# Запускаем приложение
ENTRYPOINT ["/product_accounting_service"]