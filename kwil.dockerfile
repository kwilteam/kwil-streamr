FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app

# let's add delve even if it's not debugging to make it available
RUN apk add --no-cache git \
    && go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

ARG DEBUG_PORT
ENV DEBUG_PORT=$DEBUG_PORT

# this file will help us to set the conditional env variables
RUN touch /etc/environment

# if DEBUG_PORT is set, we set GO_GCFLAGS and omit GO_LDFLAGS to avoid optimizations
RUN if [ "$DEBUG_PORT" != "" ]; then \
    echo "export GO_GCFLAGS=\"all=-N -l\"" >> /etc/environment; \
    echo "export GO_LDFLAGS=\" \"" >> /etc/environment; \
fi

RUN . /etc/environment && sh ./scripts/binary

ENTRYPOINT if [ "$DEBUG_PORT" != "" ]; then \
    echo "Running in debug mode"; \
    dlv --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient exec .build/kwil-streamr; \
else \
    echo "Running in normal mode"; \
    .build/kwil-streamr; \
fi

#CMD [".build/kwil-streamr", "--root-dir", "/root/.kwild"]