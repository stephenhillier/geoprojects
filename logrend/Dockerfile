FROM node:8.14.0-alpine

WORKDIR /app
ENV NODE_ENV production

ADD ./logrend/package*.json ./
RUN npm install

RUN mkdir -p /plotsvc/proto/plotsvc
ADD ./plotsvc/proto/plotsvc/plotsvc.proto /plotsvc/proto/plotsvc/  
ADD ./logrend/ ./

EXPOSE 8081

CMD npm run start
