# --- Builder ---
FROM golang:1.24-alpine AS builder

# Install git dan Certificate (untuk HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ojek-backend ./cmd/api/main.go

# --- Runner ---
FROM alpine:latest

# Timezone Data (WIB) & CA Certs
RUN apk --no-cache add ca-certificates tzdata

# Set Waktu ke Jakarta (WIB)
ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/ojek-backend .

EXPOSE 8080

CMD ["./ojek-backend"]
