# Etapa 1: Construcción del proyecto Next.js
FROM node:18-alpine

WORKDIR /app

# Instalamos las dependencias
COPY package*.json ./
RUN npm install

# Copiamos el código fuente y construimos el proyecto
COPY . .

# Exponemos el puerto en el que correrá el frontend
EXPOSE 3000

# Set environment variable for development
ENV NODE_ENV=development

# Ejecutamos Next.js en modo desarrollo
CMD ["npm", "run", "dev"]
