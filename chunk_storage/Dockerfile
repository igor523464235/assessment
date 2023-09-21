FROM golang:1.21.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o myapp chunk_storage/cmd/main.go

FROM alpine:3.18.3

VOLUME /data

WORKDIR /app/

COPY --from=builder /app/myapp /app/

EXPOSE 8000

ENTRYPOINT ["/app/myapp"]