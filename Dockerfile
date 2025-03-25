FROM golang:1.23.6-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go env -w CGO_ENABLED=1
RUN go build -o grpc_server ./cmd/sso/main.go

# Второй этап (создание финального контейнера)
FROM alpine:latest
WORKDIR /root/

# Создаем директорю для конфигов
RUN mkdir config

# Создаем директорию для БД
RUN mkdir storage

# Копируем конфиги
COPY --from=builder ./app/config/docker/local.yaml ./config

# Копируем БД
COPY --from=builder ./app/storage/sso.db ./storage

# Копируем переменные окружения
COPY --from=builder ./app/.env .
COPY --from=builder ./app/example.env .

COPY --from=builder /app/grpc_server .
EXPOSE 44044

CMD ["./grpc_server"]
