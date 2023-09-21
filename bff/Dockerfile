FROM golang:1.21.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o myapp bff/cmd/main.go

FROM alpine:3.18.3

WORKDIR /app/

COPY --from=builder /app/myapp /app/

COPY ./bff/wait-for-it.sh /app/wait-for-it.sh

RUN chmod +x /app/wait-for-it.sh

ENTRYPOINT ["/app/myapp"]