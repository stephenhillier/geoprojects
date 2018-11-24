FROM node:9.11.1-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.15.6-alpine
COPY --from=build-stage /app/dist /usr/share/nginx/html
WORKDIR /etc/nginx/sites-enabled/default
COPY default.conf /etc/nginx/sites-enabled/default/
EXPOSE 80
CMD ["nginx", "-c", "/etc/nginx/nginx.conf", "-g", "daemon off;"]
