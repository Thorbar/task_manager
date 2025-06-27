# Frontend Dockerfile (frontend-svelte/Dockerfile)

FROM node:20-alpine

WORKDIR /app

# Copy files needed
COPY . .

# Install dependencies
RUN npm install

# Expose dev port server
EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]
