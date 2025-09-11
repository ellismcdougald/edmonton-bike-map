# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ebm-backend ./cmd/backend

# Final
FROM alpine:latest
WORKDIR /app
COPY ebm-backend .
RUN chmod +x ./ebm-backend
EXPOSE 8000
CMD ["./ebm-backend"]
