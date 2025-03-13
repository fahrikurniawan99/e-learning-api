# Stage 1: Build stage (menggunakan image Golang)
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum lalu download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build aplikasi
RUN go build -o main .

# Stage 2: Run stage (menggunakan image yang lebih ringan)
FROM alpine:latest

# Install dependencies (jika perlu)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Expose port (ganti sesuai dengan port aplikasi Go)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
