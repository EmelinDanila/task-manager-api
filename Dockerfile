# Используем официальный образ Go для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение
RUN go build -o task-manager-api .

# Финальный контейнер
FROM alpine:latest

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/task-manager-api .

# Загружаем переменные окружения
ENV GO_ENV=production

# Указываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./task-manager-api"]
