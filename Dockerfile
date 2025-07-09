
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./... ./Api/...

WORKDIR /app/Api/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/

COPY --from=builder /app/Api/cmd/server .
COPY .env .

EXPOSE 8080
CMD ["./server"]
