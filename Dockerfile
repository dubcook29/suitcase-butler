FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

WORKDIR /app/service/wmpci
RUN go mod tidy

WORKDIR /app/service/built_wmp
RUN go mod tidy

WORKDIR /app/service
RUN go mod tidy
RUN GOARCH=amd64 go build -o service_runder . 

FROM node:22.15 AS frontend

WORKDIR /app/web
COPY ./web/package.json ./web/package-lock.json ./
RUN npm install

COPY ./web ./
RUN npm run build


FROM nginx:alpine

WORKDIR /app
COPY --from=builder /app/service/service_runder ./service/service_runder
COPY --from=frontend /app/web/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose Nginx service port
EXPOSE 80

ENV DB_ADDRESS=172.17.0.3:27017
ENV DB_USERNAME=
ENV DB_PASSWORD=

CMD ["sh", "-c", "/app/service/service_runder -dbs $DB_ADDRESS -user $DB_USERNAME -pass $DB_PASSWORD"]