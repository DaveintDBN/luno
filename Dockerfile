# Multi-stage build for Luno Trading Bot
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/bot/main.go

FROM alpine:latest
WORKDIR /root/
# Copy binary and web assets
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/config/config.json ./config/config.json
EXPOSE 8080
ENTRYPOINT ["/bin/sh", "-c", "./main --api_key_id $API_KEY_ID --api_key_secret $API_KEY_SECRET --config config/config.json"]
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/healthz || exit 1
