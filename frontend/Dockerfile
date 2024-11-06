# Etapa 1: Construcción del proyecto Next.js
FROM node:18-alpine AS builder

WORKDIR /app

# Instalamos las dependencias
COPY package*.json ./
RUN npm install

# Copiamos el código fuente y construimos el proyecto
COPY . ./
RUN npm run build

# Etapa 2: Imagen ligera para producción
FROM node:18-alpine

WORKDIR /app

# Copiamos los archivos del build generado en la etapa anterior
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# Exponemos el puerto en el que correrá el frontend
EXPOSE 3000

# Ejecutamos Next.js en modo producción
CMD ["npm", "run", "start"]
