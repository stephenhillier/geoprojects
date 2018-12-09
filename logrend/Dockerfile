FROM node:8.14.0-alpine

WORKDIR /app
ENV NODE_ENV production

ADD ./package*.json ./
RUN npm install

ADD . .

EXPOSE 8081

CMD npm run start
