FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app

# let's add delve even if it's not debugging to make it available
RUN apk add --no-cache git \
    && go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

RUN sh ./scripts/binary

RUN sh ./scripts/kwil_binaries.sh

ARG DEBUG_PORT

# if there's a debug port, we use -gcflags "all=-N -l"
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/.build/kwild /app/cmd/kwild/main.go
RUN if [ "$DEBUG_PORT" != "" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /app/.build/kwild /app/cmd/kwild/main.go; \
else \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/.build/kwild /app/cmd/kwild/main.go; \
fi

FROM busybox:1.35.0-uclibc as busybox

FROM gcr.io/distroless/static-debian12

ARG DEBUG_PORT
ENV DEBUG_PORT=$DEBUG_PORT

COPY --from=busybox /bin/sh /bin/sh

WORKDIR /app
COPY --from=build /app/.build/kwil-streamr ./kwil-streamr

COPY --from=build /app/.build/kwil-admin /app/kwil-admin

COPY --from=build /go/bin/dlv /usr/local/bin/dlv

RUN ./kwil-admin setup init --chain-id logstore-test -o /root/.kwild

EXPOSE 50051 50151 8080 26656 26657
