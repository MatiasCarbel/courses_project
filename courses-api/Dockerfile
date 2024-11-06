# Etapa 1: Construcción del binario
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copiamos los archivos de Go y descargamos las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código fuente y construimos la aplicación
COPY . .
RUN go build -o /courses-api

# Etapa 2: Imagen ligera para producción
FROM alpine:latest

WORKDIR /root/

# Copiamos el binario de la etapa de construcción
COPY --from=builder /courses-api .

# Exponemos el puerto en el que correrá el servicio
EXPOSE 8002

# Ejecutamos el binario
CMD ["./courses-api"]
