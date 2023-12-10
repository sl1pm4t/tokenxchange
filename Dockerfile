FROM golang:1.21-alpine

RUN apk add --no-cache --update alpine-sdk bash

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir bin && go build -v -o bin/tokenxchange

FROM alpine:3.19

RUN mkdir -p /app/bin
COPY --from=0 /app/bin/tokenxchange /app/bin/
