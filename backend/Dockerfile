# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ebm-backend ./cmd/server
# Install golang-migrate
RUN apk add --no-cache curl tar gzip \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /app/migrate

# Final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ebm-backend .
COPY --from=builder /app/migrate .
COPY ./migrations ./migrations
RUN chmod +x ./ebm-backend ./migrate
EXPOSE 8080
CMD ["./ebm-backend"]
