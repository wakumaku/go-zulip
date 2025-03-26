FROM caddy:2.9.1-alpine AS caddy
FROM golang:1.24.1-alpine AS dev

COPY --from=caddy /usr/bin/caddy /usr/bin/caddy

RUN apk add --no-cache \
    curl \
    git \
    gcc \
    musl-dev

# Install golang tooling
WORKDIR /tools
COPY go.mod go.sum ./
RUN go mod tidy
WORKDIR /

ENV GOFLAGS="-buildvcs=false"
