FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app

# let's add delve even if it's not debugging to make it available
RUN apk add --no-cache git \
    && go install github.com/go-delve/delve/cmd/dlv@latest

COPY ./scripts/kwil_binaries.sh ./scripts/kwil_binaries.sh

RUN sh ./scripts/kwil_binaries.sh

COPY . .

RUN sh ./scripts/binary

FROM busybox:1.35.0-uclibc as busybox

FROM gcr.io/distroless/static-debian12

ARG DEBUG_PORT
ENV DEBUG_PORT=$DEBUG_PORT

COPY --from=busybox /bin/sh /bin/sh

WORKDIR /app
COPY --from=build /app/.build/kwil-streamr ./kwil-streamr

CMD ["./kwil-streamr"]