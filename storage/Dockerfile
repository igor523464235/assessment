# Используйте официальный образ Golang
FROM golang:1.20.2 AS builder

# Установите рабочую директорию
WORKDIR /app

# Сначала копируйте go.mod и go.sum, чтобы воспользоваться кэшированием слоев Docker
COPY go.mod go.sum ./

# Загрузите зависимости (этот шаг будет кэширован, если go.mod и go.sum не изменились)
RUN go mod download

# Скопируйте исходные файлы приложения в контейнер
COPY . .

# Соберите приложение
RUN CGO_ENABLED=0 go build -o myapp storage/cmd/main.go

# Используйте отдельный этап для запуска приложения
FROM alpine:3.16

# Создайте директорию для данных
VOLUME /data

# Установите рабочую директорию
WORKDIR /app/

# Копируйте исполняемый файл из этапа builder
COPY --from=builder /app/myapp /app/

# Укажите команду для запуска приложения

# EXPOSE 8000
ENTRYPOINT ["/app/myapp"]