# build/migrator/Dockerfile

# 1) Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build the migrator binary
WORKDIR /app/cmd/migrator
RUN go build -o /app/bin/migrator

# 2) Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/migrator /app/migrator
COPY config config

ENTRYPOINT ["./migrator"]
