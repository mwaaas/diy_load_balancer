FROM golang:1.12.4-alpine3.9

RUN mkdir -p /usr/src/app \
    && apk update \
    && apk add git make

WORKDIR /usr/src/app

COPY . .
