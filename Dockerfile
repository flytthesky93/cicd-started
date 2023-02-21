# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

RUN ls
COPY config.yml ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping
RUN mkdir certs
RUN ls

EXPOSE 8100 443

CMD [ "/docker-gs-ping" ]