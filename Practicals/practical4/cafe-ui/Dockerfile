FROM node:18-alpine AS build
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
RUN npm config set registry https://registry.npmjs.org/
RUN npm config set fetch-timeout 600000
RUN npm cache clean --force
RUN npm ci || npm ci || npm install
COPY . ./
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
# Nginx's default port is 80
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]