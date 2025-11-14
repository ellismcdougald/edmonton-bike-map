# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod/go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build main backend server
RUN go build -o ebm-backend ./cmd/server

# Build update_db binary
RUN go build -o update_db ./cmd/update_db/main.go

# Install golang-migrate
RUN apk add --no-cache curl tar gzip \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /app/migrate

# Final runtime stage
FROM alpine:latest
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/ebm-backend .
COPY --from=builder /app/update_db .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/cmd/update_db/query.osmql ./query.osmql

# Make binaries executable
RUN chmod +x ./ebm-backend ./update_db ./migrate

EXPOSE 8080

# Default command is server
CMD ["./ebm-backend"]
