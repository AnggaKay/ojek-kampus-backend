# --- STAGE 1: Builder (Memasak Kode) ---
FROM golang:1.23-alpine AS builder

# Install git (kadang diperlukan untuk download modul)
RUN apk add --no-cache git

# Set folder kerja di dalam container
WORKDIR /app

# Copy file dependency dulu (biar cache-nya efisien)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi menjadi binary bernama 'main'
# CGO_ENABLED=0 membuat binary static murni (lebih ringan & kompatibel)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# --- STAGE 2: Runner (Menyajikan Aplikasi) ---
FROM alpine:latest

# Install sertifikat SSL (Penting untuk request HTTPS ke API luar/Firebase)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy hasil build (binary) dari Stage 1
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
