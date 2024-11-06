# Etapa 1: Construcción del binario
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copiamos los archivos de Go y descargamos las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código fuente y construimos la aplicación
COPY . ./
RUN go build -o /search-api ./main.go

# Etapa 2: Imagen ligera para la producción
FROM alpine:latest

WORKDIR /root/

# Copiamos el binario de la etapa de construcción
COPY --from=builder /search-api .

# Exponemos el puerto en el que correrá el servicio
EXPOSE 8003

# Ejecutamos el binario
CMD ["./search-api"]
