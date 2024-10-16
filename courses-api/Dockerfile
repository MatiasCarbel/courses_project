# Usamos una imagen base de Go
FROM golang:1.18-alpine

WORKDIR /app

# Copiar los archivos del proyecto
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compilar el binario
RUN go build -o /courses-api

# Definir el comando de ejecución
CMD ["/courses-api"]
