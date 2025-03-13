# Stage 2: Run stage (menggunakan image yang lebih ringan)
FROM alpine:latest

# Install dependencies (jika perlu)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary dari stage builder
COPY . .

# Expose port (ganti sesuai dengan port aplikasi Go)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]