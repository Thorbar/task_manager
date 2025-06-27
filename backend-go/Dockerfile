# Backend Dockerfile (backend-go/Dockerfile)

# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Necesario para compilar correctamente
RUN apk add --no-cache git

WORKDIR /app
COPY . .

# Descarga dependencias
RUN go mod tidy

# Compila desde ./cmd/main.go
RUN go build -o server ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

# Bash y curl installation
RUN apk add --no-cache bash curl

WORKDIR /app

# Copia el ejecutable y el .env
COPY --from=builder /app/server .
COPY .env .env

# Copia la carpeta assets (importante para traducciones)
COPY assets ./assets

EXPOSE 8080

# Usa APP_ENV si se define, y ejecuta el binario
CMD ["./server"]
