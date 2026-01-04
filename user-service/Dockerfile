# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app ./cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /go-app .

# Copy .env if needed (or rely on docker-compose env_file)
# COPY .env .

EXPOSE 8080
CMD ["./go-app"]