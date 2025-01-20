# Билд-образ на основе Golang
FROM golang:1.23-alpine AS builder

# Рабочая директория для приложения
WORKDIR /app

# Копируем файлы для модулей
COPY go.mod go.sum ./

# Обновляем зависимости
RUN go mod tidy

# Копируем исходный код
COPY . .

# Компилируем приложение
RUN go build -o todo-api

# Минимальный финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный файл из билдера
COPY --from=builder /app/todo-api .

# Запускаем приложение
CMD ["./todo-api"]