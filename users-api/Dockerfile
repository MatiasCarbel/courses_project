# Etapa 1: Construcción del binario
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copiamos los archivos de Go y descargamos las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código fuente y construimos la aplicación
COPY . ./
RUN go build -o /users-api ./main.go

# Etapa 2: Imagen ligera para la producción
FROM alpine:latest

WORKDIR /root/

# Instalar bash
RUN apk add --no-cache bash

# Copiamos el binario de la etapa de construcción
COPY --from=builder /users-api .

# Añadimos el script wait-for-it
COPY wait-for-it.sh .

# Hacemos que el script wait-for-it sea ejecutable
RUN chmod +x wait-for-it.sh

# Exponemos el puerto en el que correrá el servicio
EXPOSE 8001

# Usamos wait-for-it para asegurarnos que MySQL esté disponible antes de iniciar el servicio
CMD ["./wait-for-it.sh", "mysql:3306", "--", "./users-api"]
