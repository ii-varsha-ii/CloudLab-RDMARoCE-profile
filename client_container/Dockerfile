# Stage 1: Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go build -o client client.go

## Stage 2: Final stage
#FROM alpine:latest
#
#WORKDIR /app
#
#COPY --from=builder /app/client .

ENV REDIS_MASTER_IP="10.20.1.2" \
    REDIS_MASTER_PORT="6379" \
    REDIS_MASTER_PASSWORD="" \
    REDIS_MASTER_KEY="0"

ENTRYPOINT ["/app/client"]
