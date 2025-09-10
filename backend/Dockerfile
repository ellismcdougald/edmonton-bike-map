# Minimal image for your prebuilt binary
FROM alpine:latest

WORKDIR /app

# Copy the Linux/amd64 binary you built locally
COPY ebm-backend .

RUN chmod +x ./ebm-backend

EXPOSE 8000

# Run the binary
CMD ["./ebm-backend"]
