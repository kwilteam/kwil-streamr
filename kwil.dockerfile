FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN sh ./scripts/binary

CMD [".build/kwil-streamr", "--root-dir", "/root/.kwild"]