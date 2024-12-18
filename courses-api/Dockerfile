# Etapa 1: Construcción del binario
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git docker-cli

# Copiamos los archivos de Go y descargamos las dependencias
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Copiamos el código fuente y construimos la aplicación
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /courses-api

# Etapa 2: Imagen ligera para producción
FROM alpine:latest

WORKDIR /root/

# Install Docker client
RUN apk add --no-cache docker-cli

# Copiamos el binario de la etapa de construcción
COPY --from=builder /courses-api .

# Exponemos el puerto en el que correrá el servicio
EXPOSE 8002

# Ejecutamos el binario
CMD ["./courses-api"]
