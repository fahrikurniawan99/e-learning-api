# Stage 1: Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary statically linked
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Run stage
FROM alpine:latest

WORKDIR /root/

# Install dependencies
RUN apk --no-cache add ca-certificates

# Copy binary dari stage builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
