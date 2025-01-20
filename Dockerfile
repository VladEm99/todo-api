# Первый этап: сборка
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь проект и компилируем его
COPY . .
RUN go build -o todo-api

# Финальный этап: создание минимального образа
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарный файл из первого этапа
COPY --from=builder /app/todo-api .

# Указываем порт
EXPOSE 8080

# Команда для запуска приложения
CMD ["./todo-api"]