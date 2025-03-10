# Stage 1: Build
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

RUN chmod +x /root/main  # Pastikan file bisa dieksekusi

EXPOSE 8080

CMD ["./main"]
