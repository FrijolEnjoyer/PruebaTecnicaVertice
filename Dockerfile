# Multi-stage Dockerfile for PruebaTecnicaVertice

# Builder stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
WORKDIR /app/Api/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/

# Copy built binary and .env
COPY --from=builder /app/Api/cmd/server .
COPY .env .

COPY run_tests.sh .
RUN chmod +x ./run_tests.sh

EXPOSE 8080
CMD ["./server"]
