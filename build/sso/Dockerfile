# build/sso/Dockerfile

# 1) Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# grab modules
COPY go.mod go.sum ./
RUN go mod download

# copy everything
COPY . .

# build the SSO binary
WORKDIR /app/cmd/sso
RUN go build -o /app/bin/sso

# 2) Final stage
FROM alpine:3.18

# for TLS, logging, etc.
RUN apk add --no-cache ca-certificates

WORKDIR /app

# copy the binary
COPY --from=builder /app/bin/sso /app/sso

# copy config directory
COPY config config

EXPOSE 44044

ENTRYPOINT ["./sso"]
