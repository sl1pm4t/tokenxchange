FROM golang:1.21

RUN apk add --no-cache --update alpine-sdk bash

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir bin && go build -v -o bin/tokenxchange

FROM docker.io/library/ubuntu:22.04@sha256:0bced47fffa3361afa981854fcabcd4577cd43cebbb808cea2b1f33a3dd7f508

RUN mkdir -p /app/bin
COPY --from=0 /app/bin/tokenxchange /app/bin/
