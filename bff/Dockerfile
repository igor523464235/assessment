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
RUN CGO_ENABLED=0 go build -o myapp bff/cmd/main.go

# Используйте отдельный этап для запуска приложения
FROM alpine:3.16

# Установите рабочую директорию
WORKDIR /app/

# Копируйте исполняемый файл из этапа builder
COPY --from=builder /app/myapp /app/

# Копирование скрипта в образ
COPY ./bff/wait-for-it.sh /app/wait-for-it.sh

# Делаем скрипт исполняемым
RUN chmod +x /app/wait-for-it.sh

# Укажите команду для запуска приложения

EXPOSE 8001
ENTRYPOINT ["/app/myapp"]